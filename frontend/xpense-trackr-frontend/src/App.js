// src/App.js

import { Suspense, useEffect, useRef, useState } from 'react';
import {BrowserRouter as Router, Routes, Route, NavLink, Link, useLocation} from 'react-router-dom';
import './App.css';
import { faHouse, faThumbtack} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import OfflineNotice from "./components/Common/OfflineNotice";
import { navigationGroups } from './app/navigation';
import { appRoutes } from './app/routes';

function MainNavigation() {
    const { pathname } = useLocation();
    const [openGroup, setOpenGroup] = useState(null);
    const navigationRef = useRef(null);

    useEffect(() => {
        setOpenGroup(null);
    }, [pathname]);

    useEffect(() => {
        if (!openGroup) return undefined;

        const closeOnOutsidePointer = (event) => {
            const openedMenu = navigationRef.current?.querySelector('details[open]');
            if (openedMenu && !openedMenu.contains(event.target)) setOpenGroup(null);
        };

        document.addEventListener('pointerdown', closeOnOutsidePointer);
        return () => document.removeEventListener('pointerdown', closeOnOutsidePointer);
    }, [openGroup]);

    return (
        <nav ref={navigationRef} className="app-navigation" aria-label="Main navigation">
            <NavLink to="/" end className={({ isActive }) => `overview-link${isActive ? ' active-tab' : ''}`}>
                <FontAwesomeIcon icon={faHouse}/> Overview
            </NavLink>
            {navigationGroups.map((group) => {
                const isGroupActive = group.links.some(([, path]) => pathname.startsWith(path));
                return (
                    <details
                        className={`navigation-group${isGroupActive ? ' group-active' : ''}`}
                        key={group.label}
                        open={openGroup === group.label}
                    >
                        <summary onClick={(event) => {
                            event.preventDefault();
                            setOpenGroup((current) => current === group.label ? null : group.label);
                        }}>
                            {group.label}
                        </summary>
                        <div className="navigation-menu">
                            {group.links.map(([label, path]) => (
                                <NavLink
                                    key={path}
                                    to={path}
                                    className={({ isActive }) => isActive ? 'active-tab' : undefined}
                                >
                                    {label}
                                </NavLink>
                            ))}
                        </div>
                    </details>
                );
            })}
        </nav>
    );
}

function App() {
    return (
        <Router>
            <div className="xpense-trackr-main">
                <header className="app-header">
                    <Link to="/" className="app-brand">
                        <span className="brand-mark"><FontAwesomeIcon icon={faThumbtack}/></span>
                        <span>
                            <span className="brand-eyebrow">PERSONAL FINANCE</span>
                            <h1>Expense Tracker</h1>
                        </span>
                    </Link>
                    <MainNavigation/>
                </header>
                <OfflineNotice />
                <Suspense fallback={<main className="ledger-page"><p role="status">Loading page…</p></main>}>
                    <Routes>
                        {appRoutes.map(([path, Component]) => (
                            <Route key={path} path={path} element={<Component/>}/>
                        ))}
                    </Routes>
                </Suspense>
            </div>
        </Router>
    );
}

export default App;
