// ========================================
// 电商平台 - 订单和用户页面（续）
// ========================================

// 继续 editAddress 函数的模态框
async function editAddressContinue(id) {
    // 此函数已在前面完成
}

async function updateAddress() {
    const data = {
        id: parseInt($('#addr-id').value),
        receiver_name: $('#addr-name').value,
        phone: $('#addr-phone').value,
        province: $('#addr-province').value,
        city: $('#addr-city').value,
        district: $('#addr-district').value,
        detail_address: $('#addr-detail').value,
        is_default: $('#addr-default').checked ? 1 : 0
    };

    const res = await api('/address/update', {
        method: 'POST',
        body: data
    });

    if (res && res.code === 200) {
        showToast('更新成功', 'success');
        closeModal();
        renderAddresses();
    } else {
        showToast(res?.message || '更新失败', 'error');
    }
}

async function deleteAddress(id) {
    if (!confirm('确定要删除这个地址吗？')) return;
    
    const res = await api('/address/delete?id=' + id, { method: 'DELETE' });
    
    if (res && res.code === 200) {
        showToast('删除成功', 'success');
        renderAddresses();
    } else {
        showToast(res?.message || '删除失败', 'error');
    }
}

// ========================================
// 商家中心
// ========================================
async function renderSellerPage(path) {
    if (!App.user || App.user.role !== 'seller') {
        showToast('请先登录商家账号', 'warning');
        navigate('/login');
        return;
    }

    const content = $('#app-content');
    
    const sidebarHtml = `
        <div class="sidebar">
            <div class="text-center mb-3">
                <div class="user-avatar" style="width: 60px; height: 60px; font-size: 1.5rem; margin: 0 auto;">
                    ${App.user.username.charAt(0).toUpperCase()}
                </div>
                <h4 class="mt-2">${App.seller?.shop_name || '我的店铺'}</h4>
                <span class="order-status ${App.seller?.status === 1 ? 'status-completed' : 'status-pending'}">
                    ${App.seller?.status === 1 ? '已认证' : '待审核'}
                </span>
            </div>
            <ul class="sidebar-menu">
                <li class="sidebar-item">
                    <a href="/seller/dashboard" class="sidebar-link ${path === '/seller/dashboard' ? 'active' : ''}">📊 数据概览</a>
                </li>
                <li class="sidebar-item">
                    <a href="/seller/products" class="sidebar-link ${path === '/seller/products' ? 'active' : ''}">📦 商品管理</a>
                </li>
                <li class="sidebar-item">
                    <a href="/seller/orders" class="sidebar-link ${path === '/seller/orders' ? 'active' : ''}">📋 订单管理</a>
                </li>
                <li class="sidebar-item">
                    <a href="/seller/reviews" class="sidebar-link ${path === '/seller/reviews' ? 'active' : ''}">⭐ 评价管理</a>
                </li>
                <li class="sidebar-item">
                    <a href="/seller/settings" class="sidebar-link ${path === '/seller/settings' ? 'active' : ''}">⚙️ 店铺设置</a>
                </li>
            </ul>
        </div>
    `;

    if (path === '/seller/dashboard') {
        await renderSellerDashboard(content, sidebarHtml);
    } else if (path === '/seller/products') {
        await renderSellerProducts(content, sidebarHtml);
    } else if (path === '/seller/orders') {
        await renderSellerOrders(content, sidebarHtml);
    } else if (path === '/seller/reviews') {
        await renderSellerReviews(content, sidebarHtml);
    } else if (path === '/seller/settings') {
        await renderSellerSettings(content, sidebarHtml);
    } else {
        navigate('/seller/dashboard');
    }
}

async function renderSellerDashboard(content, sidebarHtml) {
    const statsRes = await api('/seller/order/stats');
    const stats = statsRes?.data || {};

    content.innerHTML = `
        <div style="display: grid; grid-template-columns: 260px 1fr; gap: 20px;">
            ${sidebarHtml}
            <div>
                <h2 class="mb-3">📊 数据概览</h2>
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
                        <h3 class="card-title">💰 销售统计</h3>
                    </div>
                    <div class="card-body">
                        <div class="stat-card" style="background: linear-gradient(135deg, var(--primary), var(--primary-dark)); color: white;">
                            <div class="stat-icon" style="background: rgba(255,255,255,0.2); color: white;">💵</div>
                            <div class="stat-content">
                                <h3 style="color: white;">${formatPrice(stats.total_sales || 0)}</h3>
                                <p style="color: rgba(255,255,255,0.8);">总销售额</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    `;
}

async function renderSellerProducts(content, sidebarHtml) {
    content.innerHTML = `
        <div style="display: grid; grid-template-columns: 260px 1fr; gap: 20px;">
            ${sidebarHtml}
            <div>
                <div class="d-flex justify-between align-center mb-3">
                    <h2>📦 商品管理</h2>
                    <button class="btn btn-primary" onclick="showAddProductModal()">+ 添加商品</button>
                </div>
                <div class="card">
                    <div class="tabs">
                        <div class="tab-item active" onclick="loadSellerProducts(-1, this)">全部</div>
                        <div class="tab-item" onclick="loadSellerProducts(0, this)">待审核</div>
                        <div class="tab-item" onclick="loadSellerProducts(1, this)">已上架</div>
                        <div class="tab-item" onclick="loadSellerProducts(2, this)">已下架</div>
                    </div>
                    <div id="seller-products-container">
                        <div class="loading"><div class="spinner"></div></div>
                    </div>
                </div>
            </div>
        </div>
    `;

    loadSellerProducts(-1);
}

async function loadSellerProducts(status, tab) {
    if (tab) {
        $$('.tab-item').forEach(el => el.classList.remove('active'));
        tab.classList.add('active');
    }

    const container = $('#seller-products-container');
    container.innerHTML = '<div class="loading"><div class="spinner"></div></div>';

    const params = new URLSearchParams({ page: 1, size: 20 });
    if (status >= 0) params.append('status', status);

    const res = await api('/seller/products?' + params.toString());
    const products = res?.data || [];

    if (products.length === 0) {
        container.innerHTML = '<div class="empty-state"><div class="empty-icon">📦</div><p class="empty-text">暂无商品</p></div>';
        return;
    }

    container.innerHTML = `
        <div class="table-wrapper">
            <table class="table">
                <thead>
                    <tr>
                        <th>商品</th>
                        <th>价格</th>
                        <th>库存</th>
                        <th>销量</th>
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
                            <td>${p.stock}</td>
                            <td>${p.sales}</td>
                            <td>
                                <span class="order-status ${p.status === 1 ? 'status-completed' : p.status === 0 ? 'status-pending' : 'status-cancelled'}">
                                    ${p.status === 1 ? '已上架' : p.status === 0 ? '待审核' : '已下架'}
                                </span>
                            </td>
                            <td>
                                <button class="btn btn-sm btn-secondary" onclick="editProduct(${p.id})">编辑</button>
                                <button class="btn btn-sm btn-danger" onclick="deleteProduct(${p.id})">删除</button>
                            </td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        </div>
    `;
}

async function showAddProductModal() {
    const catRes = await api('/categories');
    const categories = catRes?.data || [];

    showModal('添加商品', `
        <form id="product-form">
            <div class="form-group">
                <label class="form-label">商品名称 *</label>
                <input type="text" class="form-input" id="prod-name" required>
            </div>
            <div class="form-group">
                <label class="form-label">商品分类</label>
                <select class="form-input form-select" id="prod-category">
                    <option value="0">请选择分类</option>
                    ${categories.map(c => `<option value="${c.id}">${c.name}</option>`).join('')}
                </select>
            </div>
            <div class="form-group">
                <label class="form-label">价格 *</label>
                <input type="number" class="form-input" id="prod-price" step="0.01" required>
            </div>
            <div class="form-group">
                <label class="form-label">原价</label>
                <input type="number" class="form-input" id="prod-original-price" step="0.01">
            </div>
            <div class="form-group">
                <label class="form-label">库存 *</label>
                <input type="number" class="form-input" id="prod-stock" required>
            </div>
            <div class="form-group">
                <label class="form-label">商品图片URL</label>
                <input type="text" class="form-input" id="prod-images" placeholder="多个图片用逗号分隔">
            </div>
            <div class="form-group">
                <label class="form-label">商品描述</label>
                <textarea class="form-input form-textarea" id="prod-desc"></textarea>
            </div>
        </form>
    `, `
        <button class="btn btn-secondary" onclick="closeModal()">取消</button>
        <button class="btn btn-primary" onclick="saveProduct()">保存</button>
    `);
}

async function saveProduct() {
    const data = {
        name: $('#prod-name').value,
        category_id: parseInt($('#prod-category').value) || 0,
        price: parseFloat($('#prod-price').value),
        original_price: parseFloat($('#prod-original-price').value) || 0,
        stock: parseInt($('#prod-stock').value),
        images: $('#prod-images').value,
        description: $('#prod-desc').value
    };

    if (!data.name || !data.price || !data.stock) {
        showToast('请填写必填项', 'warning');
        return;
    }

    const res = await api('/seller/product/create', {
        method: 'POST',
        body: data
    });

    if (res && res.code === 200) {
        showToast('商品添加成功，等待审核', 'success');
        closeModal();
        loadSellerProducts(-1);
    } else {
        showToast(res?.message || '添加失败', 'error');
    }
}

async function editProduct(id) {
    const res = await api('/product?id=' + id);
    if (!res || res.code !== 200) {
        showToast('商品不存在', 'error');
        return;
    }

    const p = res.data;
    const catRes = await api('/categories');
    const categories = catRes?.data || [];

    showModal('编辑商品', `
        <form id="product-form">
            <input type="hidden" id="prod-id" value="${p.id}">
            <div class="form-group">
                <label class="form-label">商品名称 *</label>
                <input type="text" class="form-input" id="prod-name" value="${p.name}" required>
            </div>
            <div class="form-group">
                <label class="form-label">商品分类</label>
                <select class="form-input form-select" id="prod-category">
                    <option value="0">请选择分类</option>
                    ${categories.map(c => `<option value="${c.id}" ${c.id === p.category_id ? 'selected' : ''}>${c.name}</option>`).join('')}
                </select>
            </div>
            <div class="form-group">
                <label class="form-label">价格 *</label>
                <input type="number" class="form-input" id="prod-price" step="0.01" value="${p.price}" required>
            </div>
            <div class="form-group">
                <label class="form-label">原价</label>
                <input type="number" class="form-input" id="prod-original-price" step="0.01" value="${p.original_price || ''}">
            </div>
            <div class="form-group">
                <label class="form-label">库存 *</label>
                <input type="number" class="form-input" id="prod-stock" value="${p.stock}" required>
            </div>
            <div class="form-group">
                <label class="form-label">商品图片URL</label>
                <input type="text" class="form-input" id="prod-images" value="${p.images || ''}">
            </div>
            <div class="form-group">
                <label class="form-label">商品描述</label>
                <textarea class="form-input form-textarea" id="prod-desc">${p.description || ''}</textarea>
            </div>
            <div class="form-group">
                <label class="form-label">状态</label>
                <select class="form-input form-select" id="prod-status">
                    <option value="0" ${p.status === 0 ? 'selected' : ''}>待审核</option>
                    <option value="1" ${p.status === 1 ? 'selected' : ''}>上架</option>
                    <option value="2" ${p.status === 2 ? 'selected' : ''}>下架</option>
                </select>
            </div>
        </form>
    `, `
        <button class="btn btn-secondary" onclick="closeModal()">取消</button>
        <button class="btn btn-primary" onclick="updateProduct()">保存</button>
    `);
}

async function updateProduct() {
    const data = {
        id: parseInt($('#prod-id').value),
        name: $('#prod-name').value,
        category_id: parseInt($('#prod-category').value) || 0,
        price: parseFloat($('#prod-price').value),
        original_price: parseFloat($('#prod-original-price').value) || 0,
        stock: parseInt($('#prod-stock').value),
        images: $('#prod-images').value,
        description: $('#prod-desc').value,
        status: parseInt($('#prod-status').value)
    };

    const res = await api('/seller/product/update', {
        method: 'POST',
        body: data
    });

    if (res && res.code === 200) {
        showToast('更新成功', 'success');
        closeModal();
        loadSellerProducts(-1);
    } else {
        showToast(res?.message || '更新失败', 'error');
    }
}

async function deleteProduct(id) {
    if (!confirm('确定要删除这个商品吗？')) return;
    
    const res = await api('/seller/product/delete?id=' + id, { method: 'DELETE' });
    
    if (res && res.code === 200) {
        showToast('删除成功', 'success');
        loadSellerProducts(-1);
    } else {
        showToast(res?.message || '删除失败', 'error');
    }
}

async function renderSellerOrders(content, sidebarHtml) {
    content.innerHTML = `
        <div style="display: grid; grid-template-columns: 260px 1fr; gap: 20px;">
            ${sidebarHtml}
            <div>
                <h2 class="mb-3">📋 订单管理</h2>
                <div class="card">
                    <div class="tabs">
                        <div class="tab-item active" onclick="loadSellerOrders(-1, this)">全部</div>
                        <div class="tab-item" onclick="loadSellerOrders(1, this)">待发货</div>
                        <div class="tab-item" onclick="loadSellerOrders(2, this)">已发货</div>
                        <div class="tab-item" onclick="loadSellerOrders(3, this)">已完成</div>
                    </div>
                    <div id="seller-orders-container">
                        <div class="loading"><div class="spinner"></div></div>
                    </div>
                </div>
            </div>
        </div>
    `;

    loadSellerOrders(-1);
}

async function loadSellerOrders(status, tab) {
    if (tab) {
        $$('.tab-item').forEach(el => el.classList.remove('active'));
        tab.classList.add('active');
    }

    const container = $('#seller-orders-container');
    container.innerHTML = '<div class="loading"><div class="spinner"></div></div>';

    const params = new URLSearchParams({ page: 1, size: 20 });
    if (status >= 0) params.append('status', status);

    const res = await api('/seller/orders?' + params.toString());
    const orders = res?.data || [];

    if (orders.length === 0) {
        container.innerHTML = '<div class="empty-state"><div class="empty-icon">📋</div><p class="empty-text">暂无订单</p></div>';
        return;
    }

    container.innerHTML = orders.map(order => `
        <div class="order-card">
            <div class="order-header">
                <div class="order-info">
                    <span>订单号：${order.order_no}</span>
                    <span>买家：${order.username}</span>
                    <span>${formatDate(order.created_at)}</span>
                </div>
                <span class="order-status ${getOrderStatusClass(order.status)}">${getOrderStatusText(order.status)}</span>
            </div>
            <div class="order-body">
                ${(order.items || []).map(item => `
                    <div class="order-product">
                        <div class="order-product-image">
                            <img src="${item.product_image || 'https://via.placeholder.com/80'}" alt="">
                        </div>
                        <div class="order-product-info">
                            <h4 class="order-product-name">${item.product_name}</h4>
                            <p class="order-product-price">${formatPrice(item.price)} × ${item.quantity}</p>
                        </div>
                    </div>
                `).join('')}
                <div class="mt-2" style="background: var(--gray-50); padding: 12px; border-radius: 8px;">
                    <p><strong>收货人：</strong>${order.receiver_name} ${order.receiver_phone}</p>
                    <p><strong>地址：</strong>${order.receiver_address}</p>
                    ${order.remark ? `<p><strong>备注：</strong>${order.remark}</p>` : ''}
                </div>
            </div>
            <div class="order-footer">
                <div class="order-total">
                    实付：<strong>${formatPrice(order.pay_amount)}</strong>
                </div>
                <div class="order-actions">
                    ${order.status === 1 ? `
                        <button class="btn btn-sm btn-primary" onclick="showShipModal(${order.id})">发货</button>
                    ` : ''}
                </div>
            </div>
        </div>
    `).join('');
}

function showShipModal(orderId) {
    showModal('订单发货', `
        <div class="form-group">
            <label class="form-label">物流单号</label>
            <input type="text" class="form-input" id="tracking-no" placeholder="请输入物流单号">
        </div>
    `, `
        <button class="btn btn-secondary" onclick="closeModal()">取消</button>
        <button class="btn btn-primary" onclick="shipOrder(${orderId})">确认发货</button>
    `);
}

async function shipOrder(orderId) {
    const trackingNo = $('#tracking-no').value;

    const res = await api('/seller/order/ship', {
        method: 'POST',
        body: { order_id: orderId, tracking_no: trackingNo }
    });

    if (res && res.code === 200) {
        showToast('发货成功', 'success');
        closeModal();
        loadSellerOrders(-1);
    } else {
        showToast(res?.message || '发货失败', 'error');
    }
}

async function renderSellerReviews(content, sidebarHtml) {
    const res = await api('/seller/reviews?page=1&size=50');
    const reviews = res?.data || [];

    content.innerHTML = `
        <div style="display: grid; grid-template-columns: 260px 1fr; gap: 20px;">
            ${sidebarHtml}
            <div>
                <h2 class="mb-3">⭐ 评价管理</h2>
                <div class="card">
                    <div class="card-body">
                        ${reviews.length === 0 ? `
                            <div class="empty-state">
                                <div class="empty-icon">⭐</div>
                                <p class="empty-text">暂无评价</p>
                            </div>
                        ` : reviews.map(r => `
                            <div class="review-item">
                                <div class="review-header">
                                    <div class="review-avatar">${r.username?.charAt(0) || '?'}</div>
                                    <div>
                                        <div class="review-user">${r.username || '匿名'}</div>
                                        <div class="review-rating">${'★'.repeat(r.rating)}${'☆'.repeat(5-r.rating)}</div>
                                    </div>
                                    <span class="text-muted" style="margin-left: auto;">${r.product_name}</span>
                                </div>
                                <p class="review-content">${r.content || '用户未填写评价'}</p>
                                ${r.reply ? `
                                    <div class="review-reply">
                                        <div class="review-reply-label">我的回复：</div>
                                        <p>${r.reply}</p>
                                    </div>
                                ` : `
                                    <button class="btn btn-sm btn-outline mt-2" onclick="showReplyModal(${r.id})">回复</button>
                                `}
                            </div>
                        `).join('')}
                    </div>
                </div>
            </div>
        </div>
    `;
}

function showReplyModal(reviewId) {
    showModal('回复评价', `
        <div class="form-group">
            <label class="form-label">回复内容</label>
            <textarea class="form-input form-textarea" id="reply-content" placeholder="输入您的回复"></textarea>
        </div>
    `, `
        <button class="btn btn-secondary" onclick="closeModal()">取消</button>
        <button class="btn btn-primary" onclick="submitReply(${reviewId})">提交回复</button>
    `);
}

async function submitReply(reviewId) {
    const reply = $('#reply-content').value;
    if (!reply) {
        showToast('请输入回复内容', 'warning');
        return;
    }

    const res = await api('/seller/review/reply', {
        method: 'POST',
        body: { review_id: reviewId, reply }
    });

    if (res && res.code === 200) {
        showToast('回复成功', 'success');
        closeModal();
        renderSellerPage('/seller/reviews');
    } else {
        showToast(res?.message || '回复失败', 'error');
    }
}

async function renderSellerSettings(content, sidebarHtml) {
    content.innerHTML = `
        <div style="display: grid; grid-template-columns: 260px 1fr; gap: 20px;">
            ${sidebarHtml}
            <div class="card">
                <div class="card-header">
                    <h3 class="card-title">⚙️ 店铺设置</h3>
                </div>
                <div class="card-body">
                    <form onsubmit="updateSellerSettings(event)">
                        <div class="form-group">
                            <label class="form-label">店铺名称</label>
                            <input type="text" class="form-input" id="shop-name" value="${App.seller?.shop_name || ''}">
                        </div>
                        <div class="form-group">
                            <label class="form-label">店铺简介</label>
                            <textarea class="form-input form-textarea" id="shop-desc">${App.seller?.shop_description || ''}</textarea>
                        </div>
                        <div class="form-group">
                            <label class="form-label">店铺Logo URL</label>
                            <input type="text" class="form-input" id="shop-logo" value="${App.seller?.shop_logo || ''}">
                        </div>
                        <button type="submit" class="btn btn-primary">保存设置</button>
                    </form>
                </div>
            </div>
        </div>
    `;
}

async function updateSellerSettings(e) {
    e.preventDefault();
    
    const res = await api('/user/seller', {
        method: 'POST',
        body: {
            shop_name: $('#shop-name').value,
            shop_description: $('#shop-desc').value,
            shop_logo: $('#shop-logo').value
        }
    });

    if (res && res.code === 200) {
        showToast('保存成功', 'success');
        // 更新本地数据
        if (App.seller) {
            App.seller.shop_name = $('#shop-name').value;
            App.seller.shop_description = $('#shop-desc').value;
            App.seller.shop_logo = $('#shop-logo').value;
        }
    } else {
        showToast(res?.message || '保存失败', 'error');
    }
}
