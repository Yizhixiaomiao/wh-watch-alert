import React, { useState, useRef, useEffect } from 'react';
import { Table, Input, Spin } from 'antd';
import { SearchOutlined } from '@ant-design/icons';

/**
 * 虚拟列表组件 - 用于优化长列表性能
 * 只渲染可见区域的行，支持虚拟滚动
 */
export const VirtualList = ({
    dataSource,
    columns,
    rowKey,
    loading,
    height = 600,
    rowHeight = 54,
    onRow,
    keyword = "",
}) => {
    const [visibleData, setVisibleData] = useState([]);
    const [scrollTop, setScrollTop] = useState(0);
    const [isSearching, setIsSearching] = useState(false);
    const [searchKeyword, setSearchKeyword] = useState("");
    const [filteredData, setFilteredData] = useState([]);
    
    const containerRef = useRef(null);
    const tableRef = useRef(null);
    const searchTimerRef = useRef(null);

    // 计算可见行数
    const getVisibleCount = () => Math.ceil(height / rowHeight);
    
    // 获取可见范围的数据
    const getVisibleRange = () => {
        const start = Math.floor(scrollTop / rowHeight);
        const end = Math.min(start + getVisibleCount(), filteredData.length);
        return { start, end };
    };

    // 过滤数据
    useEffect(() => {
        if (!keyword) {
            setFilteredData(dataSource);
            setIsSearching(false);
            return;
        }

        setIsSearching(true);
        setSearchKeyword(keyword.toLowerCase());

        const filtered = dataSource.filter(item => {
            const keywordLower = keyword.toLowerCase();
            return Object.values(item).some(
                val => 
                    typeof val === 'string' ? val.toLowerCase().includes(keywordLower) : false
            );
        });

        setFilteredData(filtered);
        setIsSearching(false);
    }, [keyword, dataSource]);

    // 更新可见数据
    useEffect(() => {
        const { start, end } = getVisibleRange();
        setVisibleData(filteredData.slice(start, end));
    }, [scrollTop, filteredData, rowHeight]);

    // 处理滚动
    const handleScroll = (e) => {
        const target = e.target;
        if (target === containerRef.current) {
            setScrollTop(target.scrollTop);
        }
    };

    // 表格行渲染函数
    const rowClassName = (record, index) => {
        return index % 2 === 0 ? 'bg-white' : 'bg-gray-50';
    };

    if (loading || isSearching) {
        return (
            <div style={{ 
                    textAlign: 'center', 
                    padding: '100px 0',
                    height: height 
                }}>
                <Spin tip={isSearching ? "搜索中..." : "加载中..."} />
            </div>
        );
    }

    return (
        <div 
            ref={containerRef}
            style={{ height, overflow: 'auto' }}
            onScroll={handleScroll}
        >
            <Table
                ref={tableRef}
                columns={columns}
                dataSource={visibleData}
                rowKey={rowKey}
                rowClassName={rowClassName}
                pagination={false}
                scroll={{ y: height }}
                components={{
                    body: {
                        row: ({ record, index }) => {
                            const actualIndex = Math.floor(scrollTop / rowHeight) + index;
                            return <div style={{ height: `${rowHeight}px`, display: 'flex', alignItems: 'center' }}>
                                {onRow ? onRow(record, actualIndex) : (
                                    <span>{record[rowKey]}</span>
                                )}
                            </div>;
                        },
                    },
                }}
            />
        </div>
    );
};

/**
 * 带搜索的虚拟列表
 */
export const SearchableVirtualList = ({
    dataSource,
    columns,
    rowKey,
    loading,
    height = 600,
    rowHeight = 54,
    onRow,
}) => {
    const [keyword, setKeyword] = useState("");
    const [searchTimer, setSearchTimer] = useState(null);
    const [isSearching, setIsSearching] = useState(false);

    const handleSearch = (value) => {
        if (searchTimer) {
            clearTimeout(searchTimer);
        }
        
        setSearchTimer(setTimeout(() => {
            setKeyword(value);
            setIsSearching(false);
        }, 300));
        setIsSearching(true);
    };

    return (
        <div>
            <Input
                placeholder="输入关键词搜索"
                prefix={<SearchOutlined />}
                allowClear
                onChange={(e) => handleSearch(e.target.value)}
                style={{ marginBottom: 16 }}
            />
            <VirtualList
                dataSource={dataSource}
                columns={columns}
                rowKey={rowKey}
                loading={loading || isSearching}
                height={height}
                rowHeight={rowHeight}
                keyword={keyword}
                onRow={onRow}
            />
        </div>
    );
};

export default VirtualList;