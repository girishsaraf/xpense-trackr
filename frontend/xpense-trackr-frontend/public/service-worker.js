const CACHE_NAME = 'xpense-trackr-static-v1';
const APP_SHELL = ['/', '/manifest.json', '/icon.svg'];
const STATIC_DESTINATIONS = new Set(['script', 'style', 'image', 'font']);

self.addEventListener('install', (event) => {
    event.waitUntil(
        caches.open(CACHE_NAME)
            .then((cache) => cache.addAll(APP_SHELL))
            .then(() => self.skipWaiting()),
    );
});

self.addEventListener('activate', (event) => {
    event.waitUntil(
        caches.keys()
            .then((cacheNames) => Promise.all(
                cacheNames
                    .filter((name) => name.startsWith('xpense-trackr-static-') && name !== CACHE_NAME)
                    .map((name) => caches.delete(name)),
            ))
            .then(() => self.clients.claim()),
    );
});

self.addEventListener('fetch', (event) => {
    const request = event.request;
    const url = new URL(request.url);
    if (request.method !== 'GET' || url.origin !== self.location.origin || url.pathname.startsWith('/api/')) {
        return;
    }

    if (request.mode === 'navigate') {
        event.respondWith(
            fetch(request).catch(async () => {
                const shell = await caches.match('/');
                if (!shell) throw new Error('The local app shell is not available offline.');
                return shell;
            }),
        );
        return;
    }

    if (!STATIC_DESTINATIONS.has(request.destination)) return;
    event.respondWith(
        fetch(request).then(async (response) => {
            if (response.ok) {
                const cache = await caches.open(CACHE_NAME);
                await cache.put(request, response.clone());
            }
            return response;
        }).catch(async (error) => {
            const cached = await caches.match(request);
            if (cached) return cached;
            throw error;
        }),
    );
});
