// API 基础配置
const API_BASE = '';

// 通用请求函数
async function request(url, options = {}) {
    try {
        const response = await fetch(url, {
            ...options,
            headers: {
                'Content-Type': 'application/json',
                ...options.headers
            }
        });
        
        const data = await response.json();
        return data;
    } catch (error) {
        console.error('请求错误:', error);
        return { success: false, message: '网络错误' };
    }
}

// 获取当前用户
async function getCurrentUser() {
    const data = await request('/api/current-user');
    return data.success ? data.user : null;
}

// 更新导航栏用户状态
async function updateNavbar() {
    const user = await getCurrentUser();
    const userSection = document.querySelector('.navbar-user');
    
    if (user) {
        userSection.innerHTML = `
            <span>欢迎，${user.username}</span>
            ${user.role === 'seller' ? '<a href="/seller">商家中心</a>' : ''}
            ${user.role === 'admin' ? '<a href="/admin">管理后台</a>' : ''}
            <a href="/cart">购物车</a>
            <a href="/orders">我的订单</a>
            <a href="/profile">个人中心</a>
            <button class="btn btn-sm btn-outline" onclick="logout()">退出</button>
        `;
    } else {
        userSection.innerHTML = `
            <a href="/login" class="btn btn-sm btn-primary">登录</a>
            <a href="/register" class="btn btn-sm btn-outline">注册</a>
        `;
    }
}

// 用户登出
async function logout() {
    const data = await request('/api/logout', { method: 'POST' });
    if (data.success) {
        showToast('已退出登录', 'success');
        setTimeout(() => {
            window.location.href = '/';
        }, 1000);
    }
}

// 显示提示信息
function showToast(message, type = 'info') {
    // 创建toast元素
    const toast = document.createElement('div');
    toast.className = `alert alert-${type}`;
    toast.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        z-index: 9999;
        min-width: 250px;
        animation: slideIn 0.3s ease-out;
    `;
    toast.textContent = message;
    
    document.body.appendChild(toast);
    
    // 3秒后自动消失
    setTimeout(() => {
        toast.style.animation = 'slideOut 0.3s ease-out';
        setTimeout(() => {
            document.body.removeChild(toast);
        }, 300);
    }, 3000);
}

// 添加动画样式
const style = document.createElement('style');
style.textContent = `
    @keyframes slideIn {
        from {
            transform: translateX(400px);
            opacity: 0;
        }
        to {
            transform: translateX(0);
            opacity: 1;
        }
    }
    
    @keyframes slideOut {
        from {
            transform: translateX(0);
            opacity: 1;
        }
        to {
            transform: translateX(400px);
            opacity: 0;
        }
    }
`;
document.head.appendChild(style);

// 格式化价格
function formatPrice(price) {
    return `¥${parseFloat(price).toFixed(2)}`;
}

// 格式化日期
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    });
}

// 订单状态映射
const ORDER_STATUS = {
    'pending_payment': { text: '待支付', class: 'warning' },
    'pending_shipment': { text: '待发货', class: 'info' },
    'shipped': { text: '已发货', class: 'info' },
    'completed': { text: '已完成', class: 'success' },
    'cancelled': { text: '已取消', class: 'secondary' }
};

// 获取订单状态标签
function getOrderStatusBadge(status) {
    const statusInfo = ORDER_STATUS[status] || { text: status, class: 'secondary' };
    return `<span class="badge badge-${statusInfo.class}">${statusInfo.text}</span>`;
}

// 商品状态映射
const PRODUCT_STATUS = {
    'pending': { text: '待审核', class: 'warning' },
    'approved': { text: '已上架', class: 'success' },
    'rejected': { text: '已拒绝', class: 'danger' }
};

// 获取商品状态标签
function getProductStatusBadge(status) {
    const statusInfo = PRODUCT_STATUS[status] || { text: status, class: 'secondary' };
    return `<span class="badge badge-${statusInfo.class}">${statusInfo.text}</span>`;
}

// 页面加载完成后初始化导航栏
document.addEventListener('DOMContentLoaded', () => {
    updateNavbar();
});
