import { useState, useEffect } from 'react';
import { API } from "../registerform";


export const useOnlineStatus = (userId: string | null) => {
  const [isOnline, setIsOnline] = useState(false);
  const token = localStorage.getItem('token');

  useEffect(() => {
    if (!userId) {
      setIsOnline(false);
      return;
    }

    let isMounted = true; 

    // Initial check
    const checkStatus = async () => {
      if (!isMounted) return; 
      
      try {
        const res = await fetch(`${API}/users/${userId}/online`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        
        if (!isMounted) return; 
        
        const data = await res.json();
        setIsOnline(data.online);
      } catch (error) {
        if (!isMounted) return; 
        
        console.error('Failed to check online status:', error);
        setIsOnline(false);
      }
    };

    checkStatus();

    const interval = setInterval(checkStatus, 30000);

    return () => {
      isMounted = false; 
      clearInterval(interval);
    };
  }, [userId, token]); 

  return isOnline;
};