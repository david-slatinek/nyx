import React, {useState} from "react";
import axios from "axios";

interface Message {
    text: string;
    isUser: boolean;
}

const Dialog = () => {
    const [messages, setMessages] = useState<Message[]>([]);
    const [inputText, setInputText] = useState("");

    const handleMessageSubmit = async (event: { preventDefault: () => void; }) => {
        event.preventDefault();
        const message = {
            text: inputText,
            isUser: true,
        };
        setMessages([...messages, message]);
        setInputText("");

        try {
            const response = await axios.post("/api/chat", {message});
            const data = response.data;
            const responseMessage = {
                text: data.message,
                isUser: false,
            };
            setMessages([...messages, responseMessage]);
        } catch (error) {
            console.error(error);
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
                        {message.text}
                    </div>
                ))}
            </div>
            <form onSubmit={handleMessageSubmit} className="fixed bottom-0 left-0 right-0 flex px-4 py-2">
                <input
                    type="text"
                    value={inputText}
                    onChange={(event) => setInputText(event.target.value)}
                    className="flex-1 rounded-lg p-2 mr-2 bg-gray-200"
                />
                <button type="submit" className="bg-blue-500 text-white rounded-lg p-2">
                    Send
                </button>
            </form>
        </div>
    );
}

export default Dialog;