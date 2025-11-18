"use client";

import { useParams } from "next/navigation";
import React from "react";

const ActivationPage = () => {
  const { token } = useParams();

  const [loading, setLoading] = React.useState(false);
  const [activated, setActivated] = React.useState(false);
  const [errorMessage, setErrorMessage] = React.useState("");

  const handleClick = async () => {
    setLoading(true);
    setErrorMessage(""); // clear old errors

    try {
      const res = await fetch(
        `http://localhost:${process.env.NEXT_PUBLIC_ADDR}/api/mentor/activate/${token}`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
          },
        }
      );

      const data = await res.json(); 

      if (!res.ok) {
        // extract error from backend response
        const msg =
          data?.error?.message ||
          data?.message ||
          "Something went wrong. Please try again.";

        setErrorMessage(msg);
        return;
      }

      setActivated(true);
    } catch (err) {
      setErrorMessage("Network error. Please try again.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  // UI Rendering
  return (
    <div className="flex w-full flex-col items-center justify-center min-h-screen">
      {!activated ? (
        <>
          <button
            onClick={handleClick}
            className="h-16 w-[30%] text-2xl text-white font-bold bg-black rounded-lg cursor-pointer"
          >
            {loading ? "Activating..." : "Activate"}
          </button>

          {errorMessage && (
            <p className="mt-5 text-red-600 font-semibold">{errorMessage}</p>
          )}
        </>
      ) : (
        <span className="cursor-pointer text-2xl text-black font-bold">Activated</span>
      )}
    </div>
  );
};

export default ActivationPage;