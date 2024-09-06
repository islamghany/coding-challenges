import { useState } from "react";
import "./App.css";

function App() {
  const [shortenedLink, setShortenedLink] = useState<string>("");
  const [error, setError] = useState<string>("");
  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const url = (e.currentTarget[0] as HTMLFormElement)?.value;
    if (!url) return;
    try {
      const response = await fetch("http://localhost:8080/shorten", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ url }),
      });
      if (!response.ok) {
        throw new Error("Failed to shorten the link");
      }

      const res = await response.json();
      if (!res?.data?.shortened) {
        throw new Error("Failed to shorten the link");
      }

      setShortenedLink(res?.data?.shortened);
    } catch (error) {
      // Handle error
      if (error instanceof Error) {
        setError(error.message);
      }
    }
  };
  return (
    <div className="container">
      {shortenedLink ? (
        <div className="form-modal shortened-link">
          <p>
            Shortened Link:{" "}
            <a
              href={shortenedLink}
              target="_blank"
              rel="noopener noreferrer"
              className="link"
            >
              {shortenedLink}
            </a>
          </p>
          <button onClick={() => setShortenedLink("")}>Shorten Another</button>
        </div>
      ) : (
        <form className="form-modal" onSubmit={handleSubmit}>
          <input type="text" placeholder="Enter you link" />
          <button type="submit">Shorten</button>
          {error && <p className="error">{error}</p>}
        </form>
      )}
    </div>
  );
}

export default App;
