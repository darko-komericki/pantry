import type { QueryClient } from '@tanstack/react-query'
import { Outlet, createRootRouteWithContext } from '@tanstack/react-router'

export interface RouterContext {
  queryClient: QueryClient
}

export const rootRoute = createRootRouteWithContext<RouterContext>()({
  component: RootLayout,
})

function RootLayout() {
  return (
    <div className="app">
      <header className="app__header">
        <h1 className="app__title">Pantry</h1>
      </header>
      <main className="app__main">
        <Outlet />
      </main>
    </div>
  )
}
