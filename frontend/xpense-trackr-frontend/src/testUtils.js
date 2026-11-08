import React from 'react';
import { render } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';

// Render a component inside a router. Extra routes are useful for asserting
// that a component navigated to another page.
export function renderWithRouter(ui, { route = '/', path = '/', extraRoutes = [] } = {}) {
    return render(
        <MemoryRouter initialEntries={[route]}>
            <Routes>
                <Route path={path} element={ui} />
                {extraRoutes.map((extraRoute) => (
                    <Route key={extraRoute.path} path={extraRoute.path} element={extraRoute.element} />
                ))}
            </Routes>
        </MemoryRouter>
    );
}

export function jsonResponse(data) {
    return Promise.resolve({ data });
}
