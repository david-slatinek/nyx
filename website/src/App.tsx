import React from "react";
import {createRoot} from "react-dom/client";

import "./index.scss";
import Header from "./components/Header";
import Dialog from "./components/Dialog";
import {BrowserRouter as Router, Route, Routes} from "react-router-dom";

const App = () => (
    <>
        <div>
            <Header/>
        </div>

        <Router>
            <div>
                <Routes>
                    <Route path="/"/>
                    <Route path="/dialog" element={<Dialog/>}/>
                </Routes>
            </div>
        </Router>
    </>
);

const rootElement = document.getElementById("app")!;
const root = createRoot(rootElement);
root.render(<App/>);
