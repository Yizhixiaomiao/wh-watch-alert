import React from 'react';
import { ConfigProvider, theme, App as AntdApp } from 'antd';
import { Helmet } from 'react-helmet';
import routes from './routes';
import { useRoutes } from 'react-router-dom';
import './index.css'
import { AppContextProvider } from './context/RuleContext'

export default function App() {
    const element = useRoutes(routes);
    const title = "Wh-Ops-Alert";

    return (
        <AppContextProvider>
            <AntdApp>
                <ConfigProvider  theme={{ algorithm: theme.defaultAlgorithm }}>
                    <Helmet>
                        <title>{title}</title>
                    </Helmet>
                    {element}
                </ConfigProvider>
            </AntdApp>
        </AppContextProvider>
    );
}