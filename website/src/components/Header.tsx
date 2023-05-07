import React from 'react';

const Header = () => {
    return (
        <header className="flex items-center justify-between px-4 py-3 bg-gray-800">
            <h1 className="text-white font-bold text-xl">Nyx</h1>
            <nav>
                <ul className="flex">
                    <li className="mr-6">
                        <a className="text-white hover:text-gray-300" href="/">Home</a>
                    </li>
                    <li className="mr-6">
                        <a className="text-white hover:text-gray-300" href="/dialog">Dialog</a>
                    </li>
                    <li className="mr-6">
                        <a className="text-white hover:text-gray-300" href="#">End</a>
                    </li>
                </ul>
            </nav>
        </header>
    );
};

export default Header;