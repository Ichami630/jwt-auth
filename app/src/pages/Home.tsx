import React, { useEffect, useState } from "react"

interface Profile {
  id: number
  email: string
  name?: string
  error?: string
}

const Home: React.FC = () => {
  const [profile, setProfile] = useState<Profile | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const token = localStorage.getItem("access_token")

        if (!token) {
          setError("No token found. Please login first.")
          setTimeout(() => {
            window.location.href = "/"
          }, 1500) 
          return
        }

        const res = await fetch("http://localhost:8085/profile", {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`,
          },
          credentials: "include",
        })
        const data: Profile = await res.json()
        if (!res.ok) {
          //handle unauthorised or expired token
          if (res.status === 401 && data.error === "Token expired" ){
            setError("Session Expired. Redirecting to login...");
            localStorage.removeItem("access_token"); //clear old token
            setTimeout(() => {
                window.location.href = "/"
            }, 1500) 
            return
          }else{
            throw new Error(data.error || "Unauthorized or error fetching profile");
          }
        }
        console.log(data)
        setProfile(data)
      } catch (err: any) {
        setError(err.message)
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
            <strong>ID:</strong> {profile.id}
          </p>
          <p>
            <strong>Email:</strong> {profile.email}
          </p>
          {profile.name && (
            <p>
              <strong>Name:</strong> {profile.name}
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
