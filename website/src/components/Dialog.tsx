import React, {useEffect, useRef, useState} from "react";
import axios from "axios";

interface Message {
    text: string;
    isUser: boolean;
}

const Dialog = () => {
    const [messages, setMessages] = useState<Message[]>([]);
    const [inputText, setInputText] = useState("");
    const dialogEndRef = useRef<HTMLDivElement | null>(null);

    useEffect(() => {
        if (dialogEndRef.current !== null) {
            dialogEndRef.current.scrollIntoView({behavior: "smooth"});
        }
    }, [messages]);

    const handleMessageSubmit = async (event: { preventDefault: () => void; }) => {
        event.preventDefault();
        if (inputText.trim() === "") return;

        const message = {
            text: inputText,
            isUser: true,
        };

        messages.push(message)
        setInputText("");

        try {
            const response = await axios.post("http://localhost:8080/dialog",
                {
                    text: message.text,
                    dialogID: sessionStorage.getItem("dialogID"),
                });
            const data = response.data;
            const responseMessage = {
                text: data.answer,
                isUser: false,
            };
            setMessages([...messages, responseMessage]);
        } catch (error) {
            console.error(error);
            const responseMessage = {
                text: "Failed to fetch data from the API: " + error,
                isUser: false,
            };
            setMessages([...messages, responseMessage]);
        }
    };

    return (
        <div className={"mb-10 relative min-h-screen"}>
            <div className="h-screen flex flex-col justify-end px-4 py-8 space-y-4 overflow-y-auto">
                {messages.map((message, index) => (
                    <div
                        key={index}
                        className={`p-4 rounded-lg ${
                            message.isUser ? 'bg-gray-300 self-end' : 'bg-blue-400 self-start'
                        }`}
                    >
                        <p className="text-lg font-semibold">{message.text}</p>
                    </div>
                ))}
            </div>
            <div ref={dialogEndRef}/>
            <form onSubmit={handleMessageSubmit} className="fixed bottom-0 left-0 right-0 flex px-4 py-2">
                <input
                    type="text"
                    value={inputText}
                    onChange={(event) => setInputText(event.target.value)}
                    placeholder="Type your message here"
                    className="w-full mr-2 py-2 px-4 border border-gray-400 rounded-lg shadow-sm focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
                />
                <button type="submit" className="bg-blue-500 text-white rounded-lg p-2">
                    Send
                </button>
            </form>
        </div>
    );
}

export default Dialog;