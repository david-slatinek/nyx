import React from "react";
import axios from "axios";

const End = (props: { url: string; }) => {
    try {
        axios.post(props.url + "/end", {
            text: "end",
            dialogID: sessionStorage.getItem("dialogID"),
        });
    } catch (error) {
        console.error(error);
    }

    return (
        <div className="flex justify-center">
            <h1 className="text-4xl font-bold text-center text-blue-400 mt-5 w-1/5">Dialog ended!</h1>
        </div>
    );
};

export default End;