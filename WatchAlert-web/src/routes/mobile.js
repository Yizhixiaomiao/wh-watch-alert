// 移动端路由配置

import React from 'react'
import { createBrowserRouter } from 'react-router-dom'
import { MobileTicketCreate } from '../pages/mobile/ticket/create'
import { MobileTicketQuery } from '../pages/mobile/ticket/query'

export const mobileRouter = createBrowserRouter([
    {
        path: '/mobile',
        children: [
            {
                path: 'ticket/create',
                element: <MobileTicketCreate />,
            },
            {
                path: 'ticket/query',
                element: <MobileTicketQuery />,
            },
        ]
    }
])

// 如果需要集成到现有路由中，可以这样配置：
export const mobileRoutes = [
    {
        path: '/mobile/ticket/create',
        component: MobileTicketCreate,
    },
    {
        path: '/mobile/ticket/query', 
        component: MobileTicketQuery,
    }
]