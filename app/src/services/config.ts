const token = localStorage.getItem("access_token");
export const API_CONFIG = {
    BASE_URL: "http://localhost:8085",
    headers: {
        accept: "application/json",
        Authorization: `Bearer ${token}`
    }
}