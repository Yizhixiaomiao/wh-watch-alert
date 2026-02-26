import {message} from "antd";

export const copyToClipboard = (text) => {
    if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(text).then(
            () => {
                message.success('ID 已复制到剪贴板');
            },
            () => {
                fallbackCopyToClipboard(text);
            }
        );
    } else {
        fallbackCopyToClipboard(text);
    }
};

const fallbackCopyToClipboard = (text) => {
    const textArea = document.createElement("textarea");
    textArea.value = text;
    textArea.style.position = "fixed";
    textArea.style.left = "-9999px";
    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();
    
    try {
        const successful = document.execCommand('copy');
        if (successful) {
            message.success('ID 已复制到剪贴板');
        } else {
            message.error('复制失败');
        }
    } catch (err) {
        message.error('复制失败');
    }
    
    document.body.removeChild(textArea);
};