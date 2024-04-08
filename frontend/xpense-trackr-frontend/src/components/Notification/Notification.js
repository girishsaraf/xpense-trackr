import React, { useState, useEffect } from 'react';
import './Notification.css';

const Notification = ({ message, isError }) => {
    const [isVisible, setIsVisible] = useState(false);

    useEffect(() => {
        if (message) {
            setIsVisible(true);
            const timer = setTimeout(() => {
                setIsVisible(false);
            }, 3000);

            return () => clearTimeout(timer);
        }
    }, [message]);

    const notificationClass = isError ? 'error' : 'success';

    return (
        <div className={`notification ${notificationClass} ${isVisible ? 'show' : ''}`}>
            {message}
        </div>
    );
};

export default Notification;
