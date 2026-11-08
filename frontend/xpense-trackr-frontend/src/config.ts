declare const process: {
    env: {
        REACT_APP_API_URL?: string;
        NODE_ENV?: string;
    };
};

const backendUrl = process.env.REACT_APP_API_URL
    ?? (process.env.NODE_ENV === 'development' ? 'http://localhost:8080' : '');
export default backendUrl;
