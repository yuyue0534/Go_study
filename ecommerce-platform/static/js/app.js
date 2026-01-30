// ========================================
// 电商平台 - 主要JavaScript
// ========================================

// 全局状态
const App = {
    user: null,
    seller: null,
    token: null,
    cartCount: 0
};

// API基础URL
const API_BASE = '/api';

// ========================================
// 工具函数
// ========================================
function $(selector) {
    return document.querySelector(selector);
}

function $$(selector) {
    return document.querySelectorAll(selector);
}

function createElement(tag, className, innerHTML) {
    const el = document.createElement(tag);
    if (className) el.className = className;
    if (innerHTML) el.innerHTML = innerHTML;
    return el;
}

function formatPrice(price) {
    return '¥' + parseFloat(price).toFixed(2);
}

function formatDate(dateStr) {
    const date = new Date(dateStr);
    return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    });
}

function getFirstImage(images) {
    if (!images) return 'https://via.placeholder.com/400x400?text=No+Image';
    return images.split(',')[0];
}

function getOrderStatusText(status) {
    const statusMap = {
        0: '待支付',
        1: '待发货',
        2: '已发货',
        3: '已完成',
        4: '已取消',
        5: '退款中',
        6: '已退款'
    };
    return statusMap[status] || '未知';
}

function getOrderStatusClass(status) {
    const classMap = {
        0: 'status-pending',
        1: 'status-paid',
        2: 'status-shipped',
        3: 'status-completed',
        4: 'status-cancelled',
        5: 'status-pending',
        6: 'status-cancelled'
    };
    return classMap[status] || '';
}

// ========================================
// API请求
// ========================================
async function api(endpoint, options = {}) {
    const url = API_BASE + endpoint;
    const config = {
        headers: {
            'Content-Type': 'application/json',
        },
        ...options
    };

    if (App.token) {
        config.headers['Authorization'] = 'Bearer ' + App.token;
    }

    if (config.body && typeof config.body === 'object') {
        config.body = JSON.stringify(config.body);
    }

    try {
        const response = await fetch(url, config);
        const data = await response.json();
        
        if (data.code === 401) {
            // 未授权，清除登录状态
            logout();
            showToast('请先登录', 'warning');
            return null;
        }
        
        return data;
    } catch (error) {
        console.error('API Error:', error);
        showToast('网络请求失败', 'error');
        return null;
    }
}

// ========================================
// 提示消息
// ========================================
function showToast(message, type = 'info') {
    let container = $('.toast-container');
    if (!container) {
        container = createElement('div', 'toast-container');
        document.body.appendChild(container);
    }

    const toast = createElement('div', `toast toast-${type}`, `
        <span>${type === 'success' ? '✓' : type === 'error' ? '✕' : 'ℹ'}</span>
        <span>${message}</span>
    `);

    container.appendChild(toast);

    setTimeout(() => {
        toast.style.animation = 'slideIn 0.3s ease reverse';
        setTimeout(() => toast.remove(), 300);
    }, 3000);
}

// ========================================
// 模态框
// ========================================
function showModal(title, content, footer = '') {
    let overlay = $('.modal-overlay');
    if (overlay) overlay.remove();

    overlay = createElement('div', 'modal-overlay', `
        <div class="modal">
            <div class="modal-header">
                <h3 class="modal-title">${title}</h3>
                <button class="modal-close" onclick="closeModal()">✕</button>
            </div>
            <div class="modal-body">${content}</div>
            ${footer ? `<div class="modal-footer">${footer}</div>` : ''}
        </div>
    `);

    document.body.appendChild(overlay);
    setTimeout(() => overlay.classList.add('show'), 10);

    overlay.addEventListener('click', (e) => {
        if (e.target === overlay) closeModal();
    });
}

function closeModal() {
    const overlay = $('.modal-overlay');
    if (overlay) {
        overlay.classList.remove('show');
        setTimeout(() => overlay.remove(), 300);
    }
}

// ========================================
// 认证相关
// ========================================
async function checkAuth() {
    const token = localStorage.getItem('token');
    if (!token) return false;

    App.token = token;
    const res = await api('/auth/user');
    
    if (res && res.code === 200) {
        App.user = res.data.user;
        App.seller = res.data.seller;
        updateHeader();
        return true;
    }
    
    localStorage.removeItem('token');
    App.token = null;
    return false;
}

function logout() {
    localStorage.removeItem('token');
    App.user = null;
    App.seller = null;
    App.token = null;
    updateHeader();
    if (window.location.pathname !== '/' && 
        window.location.pathname !== '/login' && 
        window.location.pathname !== '/register') {
        navigate('/');
    }
}

function updateHeader() {
    const navLinks = $('.nav-links');
    if (!navLinks) return;

    if (App.user) {
        const initial = App.user.username.charAt(0).toUpperCase();
        let menuItems = `
            <a href="/user/profile" class="dropdown-item">👤 个人中心</a>
            <a href="/orders" class="dropdown-item">📦 我的订单</a>
            <a href="/user/addresses" class="dropdown-item">📍 收货地址</a>
        `;

        if (App.user.role === 'seller') {
            menuItems += `
                <div class="dropdown-divider"></div>
                <a href="/seller/dashboard" class="dropdown-item">🏪 商家中心</a>
                <a href="/seller/products" class="dropdown-item">📦 商品管理</a>
                <a href="/seller/orders" class="dropdown-item">📋 订单管理</a>
            `;
        }

        if (App.user.role === 'admin') {
            menuItems += `
                <div class="dropdown-divider"></div>
                <a href="/admin/dashboard" class="dropdown-item">⚙️ 管理后台</a>
            `;
        }

        menuItems += `
            <div class="dropdown-divider"></div>
            <a href="#" onclick="logout(); return false;" class="dropdown-item">🚪 退出登录</a>
        `;

        navLinks.innerHTML = `
            <a href="/cart" class="nav-link">
                🛒 购物车
                <span class="badge" id="cart-badge" style="display:none">0</span>
            </a>
            <div class="user-menu">
                <div class="user-avatar">${initial}</div>
                <div class="dropdown-menu">
                    ${menuItems}
                </div>
            </div>
        `;

        updateCartCount();
    } else {
        navLinks.innerHTML = `
            <a href="/login" class="nav-link">登录</a>
            <a href="/register" class="btn btn-primary btn-sm">注册</a>
        `;
    }
}

async function updateCartCount() {
    if (!App.user) return;
    
    const res = await api('/cart/count');
    if (res && res.code === 200) {
        App.cartCount = res.data.count;
        const badge = $('#cart-badge');
        if (badge) {
            badge.textContent = App.cartCount;
            badge.style.display = App.cartCount > 0 ? 'block' : 'none';
        }
    }
}

// ========================================
// 路由
// ========================================
function navigate(path) {
    history.pushState(null, '', path);
    handleRoute();
}

async function handleRoute() {
    const path = window.location.pathname;
    const content = $('#app-content');
    
    // 显示加载
    content.innerHTML = '<div class="loading"><div class="spinner"></div></div>';

    try {
        if (path === '/' || path === '/index.html') {
            await renderHome();
        } else if (path === '/login') {
            renderLogin();
        } else if (path === '/register') {
            renderRegister();
        } else if (path.startsWith('/product/')) {
            const id = path.split('/')[2];
            await renderProductDetail(id);
        } else if (path === '/cart') {
            await renderCart();
        } else if (path === '/checkout') {
            await renderCheckout();
        } else if (path === '/orders') {
            await renderOrders();
        } else if (path.startsWith('/order/')) {
            const id = path.split('/')[2];
            await renderOrderDetail(id);
        } else if (path === '/user/profile') {
            await renderUserProfile();
        } else if (path === '/user/addresses') {
            await renderAddresses();
        } else if (path.startsWith('/seller/')) {
            await renderSellerPage(path);
        } else if (path.startsWith('/admin/')) {
            await renderAdminPage(path);
        } else {
            content.innerHTML = '<div class="empty-state"><div class="empty-icon">404</div><p class="empty-text">页面不存在</p></div>';
        }
    } catch (error) {
        console.error('Route Error:', error);
        content.innerHTML = '<div class="empty-state"><div class="empty-icon">⚠️</div><p class="empty-text">加载失败</p></div>';
    }
}

// ========================================
// 页面渲染 - 首页
// ========================================
async function renderHome() {
    const content = $('#app-content');
    
    // 获取数据
    const [categoriesRes, hotRes, newRes] = await Promise.all([
        api('/categories'),
        api('/products/hot?limit=8'),
        api('/products/new?limit=8')
    ]);

    const categories = categoriesRes?.data || [];
    const hotProducts = hotRes?.data || [];
    const newProducts = newRes?.data || [];

    content.innerHTML = `
        <!-- 分类导航 -->
        <div class="category-nav">
            <div class="category-list">
                <div class="category-item active" onclick="filterByCategory(0)">
                    <span class="category-icon">🏠</span>
                    <span class="category-name">全部</span>
                </div>
                ${categories.map(cat => `
                    <div class="category-item" onclick="filterByCategory(${cat.id})">
                        <span class="category-icon">${cat.icon || '📦'}</span>
                        <span class="category-name">${cat.name}</span>
                    </div>
                `).join('')}
            </div>
        </div>

        <!-- 热销商品 -->
        <section class="mb-4">
            <div class="d-flex align-center justify-between mb-2">
                <h2>🔥 热销商品</h2>
                <a href="#" onclick="showAllProducts('sales'); return false;" class="text-primary">查看更多 →</a>
            </div>
            <div class="product-grid">
                ${hotProducts.map(renderProductCard).join('')}
            </div>
        </section>

        <!-- 新品上架 -->
        <section class="mb-4">
            <div class="d-flex align-center justify-between mb-2">
                <h2>✨ 新品上架</h2>
                <a href="#" onclick="showAllProducts('new'); return false;" class="text-primary">查看更多 →</a>
            </div>
            <div class="product-grid">
                ${newProducts.map(renderProductCard).join('')}
            </div>
        </section>

        <!-- 所有商品 -->
        <section id="all-products" class="hidden">
            <div class="d-flex align-center justify-between mb-2">
                <h2 id="products-title">所有商品</h2>
                <select class="form-input" style="width: auto;" onchange="sortProducts(this.value)">
                    <option value="">默认排序</option>
                    <option value="price_asc">价格从低到高</option>
                    <option value="price_desc">价格从高到低</option>
                    <option value="sales">销量优先</option>
                    <option value="rating">评分优先</option>
                </select>
            </div>
            <div class="product-grid" id="products-container"></div>
            <div id="pagination"></div>
        </section>
    `;
}

function renderProductCard(product) {
    const discount = product.original_price > product.price 
        ? Math.round((1 - product.price / product.original_price) * 100) 
        : 0;

    return `
        <div class="product-card" onclick="navigate('/product/${product.id}')">
            <div class="product-image">
                <img src="${getFirstImage(product.images)}" alt="${product.name}" loading="lazy">
                ${discount > 0 ? `<span class="product-badge">-${discount}%</span>` : ''}
            </div>
            <div class="product-info">
                <h3 class="product-name">${product.name}</h3>
                <div class="product-price">
                    <span class="price-current">${formatPrice(product.price)}</span>
                    ${product.original_price > product.price ? 
                        `<span class="price-original">${formatPrice(product.original_price)}</span>` : ''}
                </div>
                <div class="product-meta">
                    <div class="product-rating">
                        ${'★'.repeat(Math.floor(product.rating || 5))}
                        <span>${(product.rating || 5).toFixed(1)}</span>
                    </div>
                    <span>已售 ${product.sales || 0}</span>
                </div>
            </div>
        </div>
    `;
}

// 全局搜索和筛选状态
let currentFilters = {
    category_id: 0,
    keyword: '',
    sort: '',
    page: 1
};

async function filterByCategory(categoryId) {
    currentFilters.category_id = categoryId;
    currentFilters.page = 1;
    
    // 更新分类高亮
    $$('.category-item').forEach(el => el.classList.remove('active'));
    event.target.closest('.category-item').classList.add('active');
    
    await loadProducts();
}

async function showAllProducts(sort = '') {
    currentFilters.sort = sort;
    currentFilters.page = 1;
    await loadProducts();
}

async function sortProducts(sort) {
    currentFilters.sort = sort;
    currentFilters.page = 1;
    await loadProducts();
}

async function searchProducts() {
    const keyword = $('#search-input').value.trim();
    currentFilters.keyword = keyword;
    currentFilters.page = 1;
    await loadProducts();
}

async function loadProducts() {
    const section = $('#all-products');
    const container = $('#products-container');
    const title = $('#products-title');
    
    section.classList.remove('hidden');
    container.innerHTML = '<div class="loading"><div class="spinner"></div></div>';

    const params = new URLSearchParams({
        page: currentFilters.page,
        size: 20,
        status: 1
    });

    if (currentFilters.category_id > 0) params.append('category_id', currentFilters.category_id);
    if (currentFilters.keyword) params.append('keyword', currentFilters.keyword);
    if (currentFilters.sort) params.append('sort', currentFilters.sort);

    const res = await api('/products?' + params.toString());
    
    if (res && res.code === 200) {
        const products = res.data || [];
        
        if (currentFilters.keyword) {
            title.textContent = `搜索 "${currentFilters.keyword}" 的结果`;
        } else {
            title.textContent = '所有商品';
        }

        if (products.length === 0) {
            container.innerHTML = '<div class="empty-state"><div class="empty-icon">📦</div><p class="empty-text">暂无商品</p></div>';
        } else {
            container.innerHTML = products.map(renderProductCard).join('');
        }

        renderPagination(res.total, res.page, res.size);
    }

    // 滚动到商品区域
    section.scrollIntoView({ behavior: 'smooth' });
}

function renderPagination(total, page, size) {
    const pagination = $('#pagination');
    const totalPages = Math.ceil(total / size);
    
    if (totalPages <= 1) {
        pagination.innerHTML = '';
        return;
    }

    let html = '<div class="pagination">';
    
    html += `<div class="page-item ${page <= 1 ? 'disabled' : ''}" onclick="goToPage(${page - 1})">‹</div>`;
    
    for (let i = 1; i <= totalPages; i++) {
        if (i === 1 || i === totalPages || (i >= page - 2 && i <= page + 2)) {
            html += `<div class="page-item ${i === page ? 'active' : ''}" onclick="goToPage(${i})">${i}</div>`;
        } else if (i === page - 3 || i === page + 3) {
            html += '<div class="page-item">...</div>';
        }
    }
    
    html += `<div class="page-item ${page >= totalPages ? 'disabled' : ''}" onclick="goToPage(${page + 1})">›</div>`;
    html += '</div>';
    
    pagination.innerHTML = html;
}

async function goToPage(page) {
    if (page < 1) return;
    currentFilters.page = page;
    await loadProducts();
}
