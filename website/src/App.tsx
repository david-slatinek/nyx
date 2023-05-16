import React, {useState} from "react";

import "./index.scss";
import {createRoot} from "react-dom/client";
import {BrowserRouter as Router, Route, Routes} from "react-router-dom";
import Header from "./components/Header";
import Dialog from "./components/Dialog";
import axios from "axios";

const App = () => {
    let url = process.env.API_URL;
    if (url === undefined) {
        url = "http://localhost:8080";
    }

    window.addEventListener("beforeunload", function (event) {
        try {
            if (sessionStorage.getItem("dialogID") === null) {
                console.log("Dialog ID is null");
                return;
            }

            axios.post(url + "/end", {
                text: "end",
                dialogID: sessionStorage.getItem("dialogID"),
            }).then(r => console.log(r));
        } catch (error) {
            console.error(error);
        } finally {
            sessionStorage.removeItem("dialogID");
        }
    });

    const [error, setError] = useState("");

    fetch(url + "/dialog", {
        method: "GET",
        mode: "cors",
    })
        .then(response => response.json())
        .then(data => {
            setError("")
            sessionStorage.setItem("dialogID", data["dialogID"]);
        })
        .catch(error => {
            console.error("Error:", error);
            setError("Failed to fetch data from the API: " + error);
        });

    return (
        <>
            <div>
                <Header/>
            </div>

            <Router>
                <div>
                    {error && (
                        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative"
                             role="alert">
                            <strong className="font-bold">Error: </strong>
                            <span className="block sm:inline">{error}</span>
                            <span className="absolute top-0 bottom-0 right-0 px-4 py-3">
                            </span>
                        </div>
                    )}

                    <Routes>
                        <Route path="/" element={<Dialog url={url}/>}/>
                    </Routes>
                </div>
            </Router>
        </>
    );
};

export default App;

const rootElement = document.getElementById("app")!;
const root = createRoot(rootElement);
root.render(<App/>);
