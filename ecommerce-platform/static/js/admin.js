// ========================================
// 电商平台 - 管理员页面
// ========================================

async function renderAdminPage(path) {
    if (!App.user || App.user.role !== 'admin') {
        showToast('需要管理员权限', 'warning');
        navigate('/');
        return;
    }

    const content = $('#app-content');
    
    const sidebarHtml = `
        <div class="sidebar">
            <div class="text-center mb-3">
                <div class="user-avatar" style="width: 60px; height: 60px; font-size: 1.5rem; margin: 0 auto; background: linear-gradient(135deg, var(--danger), #dc2626);">
                    A
                </div>
                <h4 class="mt-2">管理员</h4>
            </div>
            <ul class="sidebar-menu">
                <li class="sidebar-item">
                    <a href="/admin/dashboard" class="sidebar-link ${path === '/admin/dashboard' ? 'active' : ''}">📊 数据概览</a>
                </li>
                <li class="sidebar-item">
                    <a href="/admin/users" class="sidebar-link ${path === '/admin/users' ? 'active' : ''}">👥 用户管理</a>
                </li>
                <li class="sidebar-item">
                    <a href="/admin/sellers" class="sidebar-link ${path === '/admin/sellers' ? 'active' : ''}">🏪 商家管理</a>
                </li>
                <li class="sidebar-item">
                    <a href="/admin/products" class="sidebar-link ${path === '/admin/products' ? 'active' : ''}">📦 商品管理</a>
                </li>
                <li class="sidebar-item">
                    <a href="/admin/orders" class="sidebar-link ${path === '/admin/orders' ? 'active' : ''}">📋 订单管理</a>
                </li>
                <li class="sidebar-item">
                    <a href="/admin/categories" class="sidebar-link ${path === '/admin/categories' ? 'active' : ''}">📁 分类管理</a>
                </li>
            </ul>
        </div>
    `;

    if (path === '/admin/dashboard') {
        await renderAdminDashboard(content, sidebarHtml);
    } else if (path === '/admin/users') {
        await renderAdminUsers(content, sidebarHtml);
    } else if (path === '/admin/sellers') {
        await renderAdminSellers(content, sidebarHtml);
    } else if (path === '/admin/products') {
        await renderAdminProducts(content, sidebarHtml);
    } else if (path === '/admin/orders') {
        await renderAdminOrders(content, sidebarHtml);
    } else if (path === '/admin/categories') {
        await renderAdminCategories(content, sidebarHtml);
    } else {
        navigate('/admin/dashboard');
    }
}

async function renderAdminDashboard(content, sidebarHtml) {
    const statsRes = await api('/admin/order/stats');
    const stats = statsRes?.data || {};

    content.innerHTML = `
        <div style="display: grid; grid-template-columns: 260px 1fr; gap: 20px;">
            ${sidebarHtml}
            <div>
                <h2 class="mb-3">📊 平台数据概览</h2>
                <div class="stat-grid">
                    <div class="stat-card">
                        <div class="stat-icon warning">📋</div>
                        <div class="stat-content">
                            <h3>${stats.pending || 0}</h3>
                            <p>待支付订单</p>
                        </div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-icon info">📦</div>
                        <div class="stat-content">
                            <h3>${stats.to_ship || 0}</h3>
                            <p>待发货订单</p>
                        </div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-icon primary">🚚</div>
                        <div class="stat-content">
                            <h3>${stats.shipped || 0}</h3>
                            <p>已发货订单</p>
                        </div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-icon success">✅</div>
                        <div class="stat-content">
                            <h3>${stats.completed || 0}</h3>
                            <p>已完成订单</p>
                        </div>
                    </div>
                </div>
                <div class="card">
                    <div class="card-header">
                        <h3 class="card-title">💰 平台总销售额</h3>
                    </div>
                    <div class="card-body">
                        <div class="stat-card" style="background: linear-gradient(135deg, #10b981, #059669); color: white;">
                            <div class="stat-icon" style="background: rgba(255,255,255,0.2); color: white;">💵</div>
                            <div class="stat-content">
                                <h3 style="color: white;">${formatPrice(stats.total_sales || 0)}</h3>
                                <p style="color: rgba(255,255,255,0.8);">总交易额</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    `;
}

async function renderAdminUsers(content, sidebarHtml) {
    content.innerHTML = `
        <div style="display: grid; grid-template-columns: 260px 1fr; gap: 20px;">
            ${sidebarHtml}
            <div>
                <h2 class="mb-3">👥 用户管理</h2>
                <div class="card">
                    <div class="tabs">
                        <div class="tab-item active" onclick="loadAdminUsers('', this)">全部</div>
                        <div class="tab-item" onclick="loadAdminUsers('customer', this)">普通用户</div>
                        <div class="tab-item" onclick="loadAdminUsers('seller', this)">商家</div>
                        <div class="tab-item" onclick="loadAdminUsers('admin', this)">管理员</div>
                    </div>
                    <div id="admin-users-container">
                        <div class="loading"><div class="spinner"></div></div>
                    </div>
                </div>
            </div>
        </div>
    `;

    loadAdminUsers('');
}

async function loadAdminUsers(role, tab) {
    if (tab) {
        $$('.tab-item').forEach(el => el.classList.remove('active'));
        tab.classList.add('active');
    }

    const container = $('#admin-users-container');
    container.innerHTML = '<div class="loading"><div class="spinner"></div></div>';

    const params = new URLSearchParams({ page: 1, size: 50 });
    if (role) params.append('role', role);

    const res = await api('/admin/users?' + params.toString());
    const users = res?.data || [];

    container.innerHTML = `
        <div class="table-wrapper">
            <table class="table">
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>用户名</th>
                        <th>邮箱</th>
                        <th>角色</th>
                        <th>状态</th>
                        <th>注册时间</th>
                        <th>操作</th>
                    </tr>
                </thead>
                <tbody>
                    ${users.map(u => `
                        <tr>
                            <td>${u.id}</td>
                            <td>${u.username}</td>
                            <td>${u.email}</td>
                            <td>
                                <span class="order-status ${u.role === 'admin' ? 'status-cancelled' : u.role === 'seller' ? 'status-shipped' : 'status-completed'}">
                                    ${u.role === 'admin' ? '管理员' : u.role === 'seller' ? '商家' : '用户'}
                                </span>
                            </td>
                            <td>
                                <span class="order-status ${u.status === 1 ? 'status-completed' : 'status-cancelled'}">
                                    ${u.status === 1 ? '正常' : '禁用'}
                                </span>
                            </td>
                            <td>${formatDate(u.created_at)}</td>
                            <td>
                                ${u.status === 1 ? 
                                    `<button class="btn btn-sm btn-danger" onclick="updateUserStatus(${u.id}, 0)">禁用</button>` :
                                    `<button class="btn btn-sm btn-success" onclick="updateUserStatus(${u.id}, 1)">启用</button>`
                                }
                            </td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        </div>
    `;
}

async function updateUserStatus(userId, status) {
    const res = await api('/admin/user/status', {
        method: 'POST',
        body: { user_id: userId, status }
    });

    if (res && res.code === 200) {
        showToast('更新成功', 'success');
        loadAdminUsers('');
    } else {
        showToast(res?.message || '更新失败', 'error');
    }
}

async function renderAdminSellers(content, sidebarHtml) {
    const res = await api('/admin/sellers?page=1&size=50');
    const sellers = res?.data || [];

    content.innerHTML = `
        <div style="display: grid; grid-template-columns: 260px 1fr; gap: 20px;">
            ${sidebarHtml}
            <div>
                <h2 class="mb-3">🏪 商家管理</h2>
                <div class="card">
                    <div class="table-wrapper">
                        <table class="table">
                            <thead>
                                <tr>
                                    <th>ID</th>
                                    <th>店铺名称</th>
                                    <th>店铺简介</th>
                                    <th>状态</th>
                                    <th>入驻时间</th>
                                    <th>操作</th>
                                </tr>
                            </thead>
                            <tbody>
                                ${sellers.map(s => `
                                    <tr>
                                        <td>${s.id}</td>
                                        <td>${s.shop_name}</td>
                                        <td>${(s.shop_description || '').substring(0, 30)}...</td>
                                        <td>
                                            <span class="order-status ${s.status === 1 ? 'status-completed' : s.status === 0 ? 'status-pending' : 'status-cancelled'}">
                                                ${s.status === 1 ? '已审核' : s.status === 0 ? '待审核' : '已禁用'}
                                            </span>
                                        </td>
                                        <td>${formatDate(s.created_at)}</td>
                                        <td>
                                            ${s.status === 0 ? `
                                                <button class="btn btn-sm btn-success" onclick="updateSellerStatus(${s.id}, 1)">通过</button>
                                                <button class="btn btn-sm btn-danger" onclick="updateSellerStatus(${s.id}, 2)">拒绝</button>
                                            ` : s.status === 1 ? `
                                                <button class="btn btn-sm btn-danger" onclick="updateSellerStatus(${s.id}, 2)">禁用</button>
                                            ` : `
                                                <button class="btn btn-sm btn-success" onclick="updateSellerStatus(${s.id}, 1)">启用</button>
                                            `}
                                        </td>
                                    </tr>
                                `).join('')}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
    `;
}

async function updateSellerStatus(sellerId, status) {
    const res = await api('/admin/seller/status', {
        method: 'POST',
        body: { seller_id: sellerId, status }
    });

    if (res && res.code === 200) {
        showToast('更新成功', 'success');
        renderAdminPage('/admin/sellers');
    } else {
        showToast(res?.message || '更新失败', 'error');
    }
}

async function renderAdminProducts(content, sidebarHtml) {
    content.innerHTML = `
        <div style="display: grid; grid-template-columns: 260px 1fr; gap: 20px;">
            ${sidebarHtml}
            <div>
                <h2 class="mb-3">📦 商品管理</h2>
                <div class="card">
                    <div class="tabs">
                        <div class="tab-item active" onclick="loadAdminProducts(-1, this)">全部</div>
                        <div class="tab-item" onclick="loadAdminProducts(0, this)">待审核</div>
                        <div class="tab-item" onclick="loadAdminProducts(1, this)">已上架</div>
                        <div class="tab-item" onclick="loadAdminProducts(2, this)">已下架</div>
                    </div>
                    <div id="admin-products-container">
                        <div class="loading"><div class="spinner"></div></div>
                    </div>
                </div>
            </div>
        </div>
    `;

    loadAdminProducts(-1);
}

async function loadAdminProducts(status, tab) {
    if (tab) {
        $$('.tab-item').forEach(el => el.classList.remove('active'));
        tab.classList.add('active');
    }

    const container = $('#admin-products-container');
    container.innerHTML = '<div class="loading"><div class="spinner"></div></div>';

    const params = new URLSearchParams({ page: 1, size: 50 });
    if (status >= 0) params.append('status', status);

    const res = await api('/admin/products?' + params.toString());
    const products = res?.data || [];

    container.innerHTML = `
        <div class="table-wrapper">
            <table class="table">
                <thead>
                    <tr>
                        <th>商品</th>
                        <th>价格</th>
                        <th>店铺</th>
                        <th>状态</th>
                        <th>操作</th>
                    </tr>
                </thead>
                <tbody>
                    ${products.map(p => `
                        <tr>
                            <td>
                                <div class="d-flex align-center gap-2">
                                    <img src="${getFirstImage(p.images)}" alt="" style="width: 50px; height: 50px; object-fit: cover; border-radius: 4px;">
                                    <span>${p.name.substring(0, 30)}${p.name.length > 30 ? '...' : ''}</span>
                                </div>
                            </td>
                            <td class="text-danger">${formatPrice(p.price)}</td>
                            <td>${p.seller_name || '-'}</td>
                            <td>
                                <span class="order-status ${p.status === 1 ? 'status-completed' : p.status === 0 ? 'status-pending' : 'status-cancelled'}">
                                    ${p.status === 1 ? '已上架' : p.status === 0 ? '待审核' : '已下架'}
                                </span>
                            </td>
                            <td>
                                ${p.status === 0 ? `
                                    <button class="btn btn-sm btn-success" onclick="updateProductStatus(${p.id}, 1)">通过</button>
                                    <button class="btn btn-sm btn-danger" onclick="updateProductStatus(${p.id}, 2)">拒绝</button>
                                ` : p.status === 1 ? `
                                    <button class="btn btn-sm btn-danger" onclick="updateProductStatus(${p.id}, 2)">下架</button>
                                ` : `
                                    <button class="btn btn-sm btn-success" onclick="updateProductStatus(${p.id}, 1)">上架</button>
                                `}
                            </td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        </div>
    `;
}

async function updateProductStatus(productId, status) {
    const res = await api('/admin/product/status', {
        method: 'POST',
        body: { product_id: productId, status }
    });

    if (res && res.code === 200) {
        showToast('更新成功', 'success');
        loadAdminProducts(-1);
    } else {
        showToast(res?.message || '更新失败', 'error');
    }
}

async function renderAdminOrders(content, sidebarHtml) {
    content.innerHTML = `
        <div style="display: grid; grid-template-columns: 260px 1fr; gap: 20px;">
            ${sidebarHtml}
            <div>
                <h2 class="mb-3">📋 订单管理</h2>
                <div class="card">
                    <div class="tabs">
                        <div class="tab-item active" onclick="loadAdminOrders(-1, this)">全部</div>
                        <div class="tab-item" onclick="loadAdminOrders(0, this)">待支付</div>
                        <div class="tab-item" onclick="loadAdminOrders(1, this)">待发货</div>
                        <div class="tab-item" onclick="loadAdminOrders(2, this)">已发货</div>
                        <div class="tab-item" onclick="loadAdminOrders(3, this)">已完成</div>
                    </div>
                    <div id="admin-orders-container">
                        <div class="loading"><div class="spinner"></div></div>
                    </div>
                </div>
            </div>
        </div>
    `;

    loadAdminOrders(-1);
}

async function loadAdminOrders(status, tab) {
    if (tab) {
        $$('.tab-item').forEach(el => el.classList.remove('active'));
        tab.classList.add('active');
    }

    const container = $('#admin-orders-container');
    container.innerHTML = '<div class="loading"><div class="spinner"></div></div>';

    const params = new URLSearchParams({ page: 1, size: 50 });
    if (status >= 0) params.append('status', status);

    const res = await api('/admin/orders?' + params.toString());
    const orders = res?.data || [];

    if (orders.length === 0) {
        container.innerHTML = '<div class="empty-state"><div class="empty-icon">📋</div><p class="empty-text">暂无订单</p></div>';
        return;
    }

    container.innerHTML = `
        <div class="table-wrapper">
            <table class="table">
                <thead>
                    <tr>
                        <th>订单号</th>
                        <th>买家</th>
                        <th>店铺</th>
                        <th>金额</th>
                        <th>状态</th>
                        <th>下单时间</th>
                    </tr>
                </thead>
                <tbody>
                    ${orders.map(o => `
                        <tr>
                            <td>${o.order_no}</td>
                            <td>${o.username}</td>
                            <td>${o.seller_name || '-'}</td>
                            <td class="text-danger">${formatPrice(o.pay_amount)}</td>
                            <td>
                                <span class="order-status ${getOrderStatusClass(o.status)}">${getOrderStatusText(o.status)}</span>
                            </td>
                            <td>${formatDate(o.created_at)}</td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        </div>
    `;
}

async function renderAdminCategories(content, sidebarHtml) {
    const res = await api('/categories');
    const categories = res?.data || [];

    content.innerHTML = `
        <div style="display: grid; grid-template-columns: 260px 1fr; gap: 20px;">
            ${sidebarHtml}
            <div>
                <div class="d-flex justify-between align-center mb-3">
                    <h2>📁 分类管理</h2>
                    <button class="btn btn-primary" onclick="showAddCategoryModal()">+ 添加分类</button>
                </div>
                <div class="card">
                    <div class="table-wrapper">
                        <table class="table">
                            <thead>
                                <tr>
                                    <th>ID</th>
                                    <th>图标</th>
                                    <th>分类名称</th>
                                    <th>排序</th>
                                    <th>操作</th>
                                </tr>
                            </thead>
                            <tbody>
                                ${categories.map(c => `
                                    <tr>
                                        <td>${c.id}</td>
                                        <td style="font-size: 1.5rem;">${c.icon || '📦'}</td>
                                        <td>${c.name}</td>
                                        <td>${c.sort_order}</td>
                                        <td>
                                            <button class="btn btn-sm btn-secondary" onclick="editCategory(${c.id}, '${c.name}', '${c.icon || ''}', ${c.sort_order})">编辑</button>
                                            <button class="btn btn-sm btn-danger" onclick="deleteCategory(${c.id})">删除</button>
                                        </td>
                                    </tr>
                                `).join('')}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
    `;
}

function showAddCategoryModal() {
    showModal('添加分类', `
        <form id="category-form">
            <div class="form-group">
                <label class="form-label">分类名称</label>
                <input type="text" class="form-input" id="cat-name" required>
            </div>
            <div class="form-group">
                <label class="form-label">图标（Emoji）</label>
                <input type="text" class="form-input" id="cat-icon" placeholder="如：📱">
            </div>
            <div class="form-group">
                <label class="form-label">排序</label>
                <input type="number" class="form-input" id="cat-sort" value="0">
            </div>
        </form>
    `, `
        <button class="btn btn-secondary" onclick="closeModal()">取消</button>
        <button class="btn btn-primary" onclick="saveCategory()">保存</button>
    `);
}

async function saveCategory() {
    const data = {
        name: $('#cat-name').value,
        icon: $('#cat-icon').value,
        sort_order: parseInt($('#cat-sort').value) || 0
    };

    if (!data.name) {
        showToast('请输入分类名称', 'warning');
        return;
    }

    const res = await api('/admin/category/create', {
        method: 'POST',
        body: data
    });

    if (res && res.code === 200) {
        showToast('添加成功', 'success');
        closeModal();
        renderAdminPage('/admin/categories');
    } else {
        showToast(res?.message || '添加失败', 'error');
    }
}

function editCategory(id, name, icon, sortOrder) {
    showModal('编辑分类', `
        <form id="category-form">
            <input type="hidden" id="cat-id" value="${id}">
            <div class="form-group">
                <label class="form-label">分类名称</label>
                <input type="text" class="form-input" id="cat-name" value="${name}" required>
            </div>
            <div class="form-group">
                <label class="form-label">图标（Emoji）</label>
                <input type="text" class="form-input" id="cat-icon" value="${icon}">
            </div>
            <div class="form-group">
                <label class="form-label">排序</label>
                <input type="number" class="form-input" id="cat-sort" value="${sortOrder}">
            </div>
        </form>
    `, `
        <button class="btn btn-secondary" onclick="closeModal()">取消</button>
        <button class="btn btn-primary" onclick="updateCategory()">保存</button>
    `);
}

async function updateCategory() {
    const data = {
        id: parseInt($('#cat-id').value),
        name: $('#cat-name').value,
        icon: $('#cat-icon').value,
        sort_order: parseInt($('#cat-sort').value) || 0
    };

    const res = await api('/admin/category/update', {
        method: 'POST',
        body: data
    });

    if (res && res.code === 200) {
        showToast('更新成功', 'success');
        closeModal();
        renderAdminPage('/admin/categories');
    } else {
        showToast(res?.message || '更新失败', 'error');
    }
}

async function deleteCategory(id) {
    if (!confirm('确定要删除这个分类吗？')) return;
    
    const res = await api('/admin/category/delete?id=' + id, { method: 'DELETE' });
    
    if (res && res.code === 200) {
        showToast('删除成功', 'success');
        renderAdminPage('/admin/categories');
    } else {
        showToast(res?.message || '删除失败', 'error');
    }
}
