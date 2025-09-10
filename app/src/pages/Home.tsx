import React, { useEffect, useState } from "react"

interface Profile {
  id: number
  email: string
  name?: string
  error?: string
}

interface Token {
    accessToken: string
    refreshToken?: string
    error?: string 
}

const Home: React.FC = () => {
  const [profile, setProfile] = useState<Profile | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchProfile = async () => {
        let token = localStorage.getItem("access_token")

        if (!token) {
          setError("No token found. Please login first.")
          setTimeout(() => {
            window.location.href = "/"
          }, 1500) 
          return
        }
      try {
        const res = await fetch("http://localhost:8085/profile", {
            method: "GET",
            headers: { "Authorization": `Bearer ${token}` },
            credentials: "include"
        })

        const data = await res.json();

        if(res.status === 401 && data.error === "Token expired"){
            //access token expired, call refresh token endpoint
            const refreshRes = await fetch("http://localhost:8085/refresh",{
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include" //send refresh token cookie
            })

            if(!refreshRes.ok) throw new Error("Refresh failed");

            const refreshData: Token = await refreshRes.json();
            token = refreshData.accessToken;
            localStorage.setItem("access_token",token)

            // Retry original request
            const retryRes = await fetch("http://localhost:8085/profile", {
                method: "GET",
                headers: { "Authorization": `Bearer ${token}` },
            })

            const retryData = await retryRes.json()
            if(!refreshRes.ok && retryData.error === "Token expired"){
                localStorage.removeItem("access_token"); //clear the access token
                setError("Session EXpired. Redirecting to Login...")
                setTimeout(()=>{
                    window.location.href = "/";
                },1500);
            }
            setProfile(retryData)
            return

        }

        if (!res.ok) throw new Error("Failed to fetch profile")
        setProfile(data)
        
      } catch (error:any) {
        setError(error.message)
      }
    }

    fetchProfile()
  }, [])

  return (
    <div>
      <h2>Profile Page</h2>
      {error && <p style={{ color: "red" }}>{error}</p>}
      {profile ? (
        <div>
          <p>
            <strong>ID:</strong> {profile.id || "001"}
          </p>
          <p>
            <strong>Email:</strong> {profile.email || "test@gmail.com"}
          </p>
          {profile.name && (
            <p>
              <strong>Name:</strong> {profile.name || "ichami"}
            </p>
          )}
        </div>
      ) : (
        !error && <p>Loading...</p>
      )}
    </div>
  )
}

export default Home
