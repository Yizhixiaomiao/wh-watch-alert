import React, { useState, useRef, useEffect } from 'react';
import { Image, Spin } from 'antd';

/**
 * 图片懒加载组件
 * 当图片进入视口时才加载
 * 支持错误回退
 */
export const LazyImage = ({ src, alt, placeholder, fallback, style, ...props }) => {
    const [loaded, setLoaded] = useState(false);
    const [error, setError] = useState(false);
    const imgRef = useRef(null);
    const [inView, setInView] = useState(false);
    const observerRef = useRef(null);

    const handleLoad = () => {
        setLoaded(true);
        setError(false);
    };

    const handleError = () => {
        setError(true);
        setLoaded(false);
    };

    // 使用 Intersection Observer 实现懒加载
    useEffect(() => {
        const observer = new IntersectionObserver(
            ([entry]) => {
                if (entry.isIntersecting) {
                    setInView(true);
                } else {
                    setInView(false);
                }
            },
            {
                root: null,
                rootMargin: '50px',
                threshold: 0.1,
            }
        );

        if (imgRef.current) {
            observer.observe(imgRef.current);
            observerRef.current = observer;
        }

        return () => {
            if (observerRef.current && imgRef.current) {
                observerRef.current.disconnect();
            }
        };
    }, [src]);

    const handleImageError = () => {
        handleError();
    };

    if (error && fallback) {
        return <div style={{
            width: '100%',
            height: style?.height || '100%',
            background: '#f0f0f0',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            color: '#999',
            fontSize: '14px',
            ...style,
        }}>
            {fallback}
        </div>;
    }

    return (
        <div ref={imgRef} style={{ display: 'inline-block' }}>
            {!loaded && (
                <div style={{
                    width: '100%',
                    height: style?.height || '100%',
                    background: '#f0f0f0',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    color: '#999',
                    fontSize: '14px',
                    ...style,
                }}>
                    <Spin size="small" />
                </div>
            )}
            {loaded ? (
                <Image
                    src={src}
                    alt={alt}
                    onError={handleImageError}
                    onLoad={handleLoad}
                    preview={false}
                    style={{
                        opacity: loaded ? 1 : 0,
                        transition: 'opacity 0.3s',
                        objectFit: 'contain',
                        ...style,
                    }}
                    {...props}
                />
            ) : null}
        </div>
    );
};

export default LazyImage;