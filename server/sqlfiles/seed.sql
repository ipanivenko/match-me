BEGIN;

-- Global sequence for unique email numbers across runs
CREATE SEQUENCE IF NOT EXISTS user_seed_seq START WITH 1000;

WITH
seq AS (
  SELECT nextval('user_seed_seq')::int AS seq_n
  -- This value was :add_count, hard-coded to 50 for execution from Go
  FROM generate_series(1, 100)
),

-- Insert USERS (unique emails)
new_users AS (
  INSERT INTO users (email, password_hash)
  SELECT
    'user_' || to_char(seq_n, 'FM000000') || '@example.com' AS email,
    -- Using a placeholder hash. In production, this would be a real bcrypt hash.
    '$2a$10$3b3LkuIcx1m4q1E5mMYIfu5pAMis0H.1PSOfF55nnSJ9OZ80ZEFKS' AS password_hash
  FROM seq
  RETURNING id, email
),

-- Stable 1..batch index for deterministic names/lat-lon nudges
batch_idx AS (
  SELECT id, email, row_number() OVER (ORDER BY email) AS i
  FROM new_users
),

-- ===================================
-- DATA POOLS FOR GENERATION
-- ===================================

fi_parent_names AS (
  SELECT ARRAY[
    'Aino','Olavi','Ilona','Matti','Sanna','Jukka','Riikka','Antti','Tiina','Ville',
    'Kaisa','Mikko','Laura','Teemu','Henna','Janne','Noora','Petri','Emmi','Sari',
    'Timo','Paula','Tuomas','Marika','Simo','Outi','Pekka','Heidi','Satu','Jani',
    'Tarja','Niina','Anu','Seppo','Leena','Eeva','Arto','Katja','Eero','Kirsti'
  ] AS arr
),
fi_child_names AS (
  SELECT ARRAY[
    'Elias','Sofia','Aada','Onni','Veeti','Emma','Noel','Olivia','Leo','Iida',
    'Oskari','Helmi','Elli','Aava','Mila','Lumi','Alvar','Siiri','Eero','Vilma',
    'Niilo','Nella','Elina','Aapo','Ella','Aino','Otto','Niko','Taika','Kerttu'
  ] AS arr
),
interests_pool AS (
  SELECT ARRAY['lego','books','music','dance','football','basketball','puzzle','drawing','coding','nature'] AS arr
),
play_styles_pool AS (
  SELECT ARRAY['active','calm','creative','outdoor','team','imaginative'] AS arr
),

-- English "About" text pools (5 variations each)
parent_abouts_pool AS (
  SELECT ARRAY[
    'Friendly parent living in Finland. We enjoy outdoor activities, board games, and exploring new cafes.',
    'We are a bilingual family looking for other families to share playdates and park visits with.',
    'Just moved here! We love cooking, visiting museums, and finding the best playgrounds for our little one.',
    'Easy-going family that loves nature, biking, and weekend trips. Our kid is very energetic and curious.',
    'Busy but fun-loving household. We value creativity and kindness, and are looking for similar families to connect with.'
  ] AS arr
),
child_abouts_pool AS (
  SELECT ARRAY[
    'A curious and friendly kid who loves building blocks and listening to stories.',
    'Very energetic! Loves running, climbing, and any game that involves a ball.',
    'A calm and imaginative child. Enjoys drawing, puzzles, and playing pretend.',
    'Loves music and dancing. A bit shy at first but warms up quickly and is very caring.',
    'Our little explorer. Loves being outdoors, finding bugs, and getting muddy.'
  ] AS arr
),

-- Other Finnish locations (Helsinki is handled separately)
fi_cities_pool AS (
  SELECT ARRAY['Tampere', 'Turku', 'Oulu', 'Jyväskylä', 'Kuopio'] AS arr
),
fi_lats_pool AS (
  SELECT ARRAY[61.4978, 60.4518, 65.0121, 62.2426, 62.8924] AS arr
),
fi_lons_pool AS (
  SELECT ARRAY[23.7610, 22.2666, 25.4651, 25.7473, 27.6782] AS arr
)

-- ===================================
-- INSERT PARENT PROFILES
-- ===================================
, insert_parents AS (
  INSERT INTO parent_profiles (
    user_id, name, gender, about, languages,
    address_city, lat, lon, preferred_distance_km
  )
  SELECT
    b.id,
    pn.arr[ ((b.i - 1) % cardinality(pn.arr)) + 1 ] AS name,
    CASE WHEN (b.i % 2)=0 THEN 'female' ELSE 'male' END AS gender,
    -- Cycle through the 5 English "about" texts
    p_about.arr[ ((b.i - 1) % cardinality(p_about.arr)) + 1 ] AS about,
    -- Assign 'Finnish' and sometimes 'English'
    CASE
      WHEN (b.i % 3) = 0 THEN ARRAY['Finnish', 'English']::text[]
      ELSE ARRAY['Finnish']::text[]
    END AS languages,
    
    -- *** Location Logic: 50% Helsinki, 50% other cities ***
    CASE
      WHEN (b.i % 2) = 1 THEN 'Helsinki' -- 50% (odd numbers) go to Helsinki
      ELSE cities.arr[ ((b.i / 2 - 1) % cardinality(cities.arr)) + 1 ] -- 50% (even numbers) cycle through other cities
    END AS address_city,
    
    -- Latitude (with jitter)
    CASE
      WHEN (b.i % 2) = 1 THEN 60.1699 + (((b.i % 5) - 2) * 0.002) -- Helsinki Lat
      ELSE lats.arr[ ((b.i / 2 - 1) % cardinality(lats.arr)) + 1 ] + (((b.i % 5) - 2) * 0.002) -- Other City Lat
    END AS lat,

    -- Longitude (with jitter)
    CASE
      WHEN (b.i % 2) = 1 THEN 24.9384 + (((b.i % 7) - 3) * 0.003) -- Helsinki Lon
      ELSE lons.arr[ ((b.i / 2 - 1) % cardinality(lons.arr)) + 1 ] + (((b.i % 7) - 3) * 0.003) -- Other City Lon
    END AS lon,
    
    (10 + (b.i % 21))::int AS preferred_distance_km -- 10..30
  FROM
    batch_idx b
    CROSS JOIN fi_parent_names pn
    CROSS JOIN parent_abouts_pool p_about
    CROSS JOIN fi_cities_pool cities
    CROSS JOIN fi_lats_pool lats
    CROSS JOIN fi_lons_pool lons
  RETURNING user_id
)

-- ===================================
-- INSERT CHILD PROFILES
-- ===================================
, insert_children AS (
  INSERT INTO children (
    user_id, name, birthday, gender, about_short,
    interests, activity_level, limitations, allergies, play_styles
  )
  SELECT
    b.id,
    cn.arr[ ((b.i - 1) % cardinality(cn.arr)) + 1 ] AS child_name,
    (DATE '2018-01-01' + ((b.i * 37) % 2191)::int) AS birthday, -- 2018..2023
    CASE WHEN (b.i % 2)=0 THEN 'boy' ELSE 'girl' END AS gender,
    
    -- Cycle through the 5 English "about short" texts
    c_about.arr[ ((b.i - 1) % cardinality(c_about.arr)) + 1 ] AS about_short,

    -- Get 2 unique interests
    ARRAY[
      ip.arr[ ((b.i    ) % cardinality(ip.arr)) + 1 ],
      ip.arr[ ((b.i + 3) % cardinality(ip.arr)) + 1 ]
    ]::text[] AS interests,
    
    CASE (b.i % 3) WHEN 0 THEN 'high' WHEN 1 THEN 'medium' ELSE 'low' END AS activity_level,
    
    -- Limitations (can NOT be empty)
    CASE (b.i % 6)
      WHEN 0 THEN ARRAY['gluten_free']::text[]
      WHEN 1 THEN ARRAY['lactose_intolerant']::text[]
      WHEN 2 THEN ARRAY['peanut']::text[]
      WHEN 3 THEN ARRAY['egg']::text[]
      WHEN 4 THEN ARRAY['pollen']::text[]
      ELSE ARRAY['none']::text[] -- Changed from empty to 'none'
    END AS limitations,

    -- Allergies (set to 'none' if empty)
    CASE (b.i % 4)
      WHEN 0 THEN ARRAY['none']::text[]
      WHEN 1 THEN ARRAY['peanut']::text[]
      WHEN 2 THEN ARRAY['lactose']::text[]
      ELSE ARRAY['pollen','dust']::text[]
    END AS allergies,

    -- Get 2 unique play styles
    ARRAY[
      ps.arr[ ((b.i    ) % cardinality(ps.arr)) + 1 ],
      ps.arr[ ((b.i + 2) % cardinality(ps.arr)) + 1 ]
    ]::text[] AS play_styles
  FROM
    batch_idx b
    CROSS JOIN fi_child_names cn
    CROSS JOIN interests_pool ip
    CROSS JOIN play_styles_pool ps
    CROSS JOIN child_abouts_pool c_about
  RETURNING user_id
)

-- Final SELECT to show the command completed, otherwise INSERTs return nothing
SELECT count(*) AS users_created FROM new_users;

COMMIT;

