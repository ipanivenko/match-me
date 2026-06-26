package services

import (
	"context"
	"matchme-server/database"
	"matchme-server/internal"
	"matchme-server/structs"
	"math"
	"sort"
	"strings"
	"time"
)

// Calculates scores for all candidates, according to user
// Returns sorted list of matches with compatibility scores
func CalculateMatchingScores(
	currentProfile *database.MatchingProfile,
	currentPrefs *structs.PreferencesInput,
	candidates []database.MatchingProfile,
	ctx context.Context,
) []structs.MatchScore {

	var matches []structs.MatchScore

	// Calculate compatibility scores for each candidate
	for _, candidate := range candidates {
		// Get candidate's preferences for mutual scoring
		candidatePrefs, err := database.GetUserMatchingPreferences(ctx, internal.DB, candidate.UserID)
		if err != nil {
			continue // Skip candidates with errors
		}

		// Calculate score from current user to candidate
		score1 := CalculateCompatibilityScore(currentProfile, &candidate, currentPrefs)

		// Calculate score from candidate to current user (mutual compatibility)
		score2 := CalculateCompatibilityScore(&candidate, currentProfile, candidatePrefs)

		// Final score is average of both directions for mutual compatibility
		finalScore := (score1 + score2) / 2.0

		matches = append(matches, structs.MatchScore{
			UserID:      candidate.UserID,
			Score:       finalScore,
			MutualScore: score2,
		})
	}

	// Sort by compatibility score (highest first)
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Score > matches[j].Score
	})

	return matches
}

// CalculateCompatibilityScore calculates how compatible two profiles are
// Uses weighted scoring based on user preferences
func CalculateCompatibilityScore(
	profile1, profile2 *database.MatchingProfile,
	prefs *structs.PreferencesInput,
) float64 {
	var totalScore float64  // Sum of weighted scores
	var totalWeight float64 // Sum of all weights used

	// 1. INTERESTS COMPATIBILITY
	// Compare children's interests using Jaccard similarity
	if prefs.InterestsWeight > 0 {
		score := CalculateArrayOverlap(profile1.Interests, profile2.Interests)
		totalScore += score * float64(prefs.InterestsWeight)
		totalWeight += float64(prefs.InterestsWeight)
	}

	// 2. ACTIVITY LEVEL COMPATIBILITY
	if prefs.ActivityLevelWeight > 0 {
		score := CalculateActivityCompatibility(profile1.ActivityLevel, profile2.ActivityLevel)
		totalScore += score * float64(prefs.ActivityLevelWeight)
		totalWeight += float64(prefs.ActivityLevelWeight)
	}

	// 3. LIMITATIONS COMPATIBILITY
	if prefs.LimitationsWeight > 0 {
		score := CalculateArrayOverlap(profile1.Limitations, profile2.Limitations)
		totalScore += score * float64(prefs.LimitationsWeight)
		totalWeight += float64(prefs.LimitationsWeight)
	}

	// 4. ALLERGIES COMPATIBILITY
	if prefs.AllergiesWeight > 0 {
		score := CalculateAllergiesCompatibility(profile1.Allergies, profile2.Allergies)
		totalScore += score * float64(prefs.AllergiesWeight)
		totalWeight += float64(prefs.AllergiesWeight)
	}

	// 5. PLAY STYLES COMPATIBILITY
	if prefs.PlayStylesWeight > 0 {
		score := CalculateArrayOverlap(profile1.PlayStyles, profile2.PlayStyles)
		totalScore += score * float64(prefs.PlayStylesWeight)
		totalWeight += float64(prefs.PlayStylesWeight)
	}

	// 6. AGE COMPATIBILITY (ALWAYS CONSIDERED)
	ageScore := CalculateAgeCompatibility(profile1.ChildBirthday, profile2.ChildBirthday, prefs.MaxAgeDifference)
	totalScore += ageScore * 2.0
	totalWeight += 2.0


	// FALLBACK: If all weights are 0, use basic compatibility formula
	if totalWeight == 0 {
		return CalculateBasicCompatibility(profile1, profile2)
	}

	// Normalize final score to 0-1 range
	finalScore := totalScore / totalWeight
	return math.Max(0, math.Min(1, finalScore))
}

// Individual functions for calculating compatibility between specific attributes

// CalculateArrayOverlap calculates similarity between two string arrays
// Uses Jaccard similarity coefficient: intersection / union
// Returns 1.0 for identical arrays, 0.0 for no overlap
func CalculateArrayOverlap(arr1, arr2 []string) float64 {
	if len(arr1) == 0 && len(arr2) == 0 {
		return 1.0 // Both empty arrays are considered identical
	}
	if len(arr1) == 0 || len(arr2) == 0 {
		return 0.0 // One empty, one not - no similarity
	}

	// Create set from first array for fast lookup
	set1 := make(map[string]bool)
	for _, item := range arr1 {
		set1[strings.ToLower(item)] = true
	}

	// Count overlapping items
	overlap := 0
	for _, item := range arr2 {
		if set1[strings.ToLower(item)] {
			overlap++
		}
	}

	// Jaccard similarity: intersection / union
	union := len(arr1) + len(arr2) - overlap
	if union == 0 {
		return 1.0
	}
	return float64(overlap) / float64(union)
}

// CalculateActivityCompatibility compares activity levels
// Maps low/medium/high to numbers and calculates compatibility
// Returns higher scores for closer activity levels
func CalculateActivityCompatibility(level1, level2 string) float64 {
	if level1 == "" || level2 == "" {
		return 0.5 // Neutral score when data is missing
	}

	// Map activity levels to numbers for comparison
	levels := map[string]int{
		"low":    1,
		"medium": 2,
		"high":   3,
	}

	l1, ok1 := levels[strings.ToLower(level1)]
	l2, ok2 := levels[strings.ToLower(level2)]

	if !ok1 || !ok2 {
		return 0.5 // Unknown levels get neutral score
	}

	// Calculate compatibility based on difference
	diff := math.Abs(float64(l1 - l2))
	return math.Max(0, 1.0-diff/2.0) // Max difference is 2, so divide by 2
}

// CalculateAllergiesCompatibility handles allergy compatibility
// Similar allergies can be good (parents understand each other's challenges)
// This could be enhanced to check for conflicting allergies
func CalculateAllergiesCompatibility(allergies1, allergies2 []string) float64 {
	return CalculateArrayOverlap(allergies1, allergies2)
}

// CalculateAgeCompatibility compares children's ages
// Returns higher scores for smaller age differences, 0 - if difference exceeds user's maximum preference

func CalculateAgeCompatibility(birth1, birth2 time.Time, maxDiffYears int) float64 {
	now := time.Now()
	age1 := now.Year() - birth1.Year()
	age2 := now.Year() - birth2.Year()
	
	if now.Month() < birth1.Month() || (now.Month() == birth1.Month() && now.Day() < birth1.Day()) {
		age1--
	}
	if now.Month() < birth2.Month() || (now.Month() == birth2.Month() && now.Day() < birth2.Day()) {
		age2--
	}
	
	diffYears := math.Abs(float64(age1 - age2))
	
	if diffYears > float64(maxDiffYears) {
		return 0.0
	}
	
	return math.Max(0, 1.0 - (diffYears / float64(maxDiffYears)))
}


// CalculateBasicCompatibility provides fallback scoring when all weights are 0
// Uses equal weighting for core compatibility factors
func CalculateBasicCompatibility(profile1, profile2 *database.MatchingProfile) float64 {
	// Basic formula using equal weights for core factors
	interestsScore := CalculateArrayOverlap(profile1.Interests, profile2.Interests)
	playStylesScore := CalculateArrayOverlap(profile1.PlayStyles, profile2.PlayStyles)
	ageScore := CalculateAgeCompatibility(profile1.ChildBirthday, profile2.ChildBirthday, 2) // 2 years

	return (interestsScore + playStylesScore + ageScore ) / 4.0
}
