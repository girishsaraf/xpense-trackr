import { useEffect, useState } from 'react';
import './OfflineNotice.css';

const OfflineNotice = () => {
    const [offline, setOffline] = useState(() => typeof navigator !== 'undefined' && !navigator.onLine);

    useEffect(() => {
        const markOnline = () => setOffline(false);
        const markOffline = () => setOffline(true);
        window.addEventListener('online', markOnline);
        window.addEventListener('offline', markOffline);
        return () => {
            window.removeEventListener('online', markOnline);
            window.removeEventListener('offline', markOffline);
        };
    }, []);

    if (!offline) return null;
    return (
        <p className="offline-notice" role="status">
            You’re offline. The app shell is available, but financial data isn’t stored for offline access.
        </p>
    );
};

export default OfflineNotice;
