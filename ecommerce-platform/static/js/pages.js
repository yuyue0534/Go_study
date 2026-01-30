// ========================================
// 电商平台 - 页面渲染
// ========================================

// ========================================
// 登录页面
// ========================================
function renderLogin() {
    const content = $('#app-content');
    content.innerHTML = `
        <div class="container" style="max-width: 400px; margin-top: 60px;">
            <div class="card">
                <div class="card-body">
                    <h2 class="text-center mb-3">欢迎回来</h2>
                    <form onsubmit="handleLogin(event)">
                        <div class="form-group">
                            <label class="form-label">用户名</label>
                            <input type="text" class="form-input" id="login-username" required placeholder="请输入用户名">
                        </div>
                        <div class="form-group">
                            <label class="form-label">密码</label>
                            <input type="password" class="form-input" id="login-password" required placeholder="请输入密码">
                        </div>
                        <button type="submit" class="btn btn-primary btn-block">登录</button>
                    </form>
                    <p class="text-center mt-2 text-muted">
                        还没有账号？<a href="/register" class="text-primary">立即注册</a>
                    </p>
                </div>
            </div>
        </div>
    `;
}

async function handleLogin(e) {
    e.preventDefault();
    
    const username = $('#login-username').value;
    const password = $('#login-password').value;

    const res = await api('/auth/login', {
        method: 'POST',
        body: { username, password }
    });

    if (res && res.code === 200) {
        App.token = res.data.token;
        App.user = res.data.user;
        App.seller = res.data.seller;
        localStorage.setItem('token', res.data.token);
        
        showToast('登录成功', 'success');
        updateHeader();
        navigate('/');
    } else {
        showToast(res?.message || '登录失败', 'error');
    }
}

// ========================================
// 注册页面
// ========================================
function renderRegister() {
    const content = $('#app-content');
    content.innerHTML = `
        <div class="container" style="max-width: 400px; margin-top: 60px;">
            <div class="card">
                <div class="card-body">
                    <h2 class="text-center mb-3">创建账号</h2>
                    <form onsubmit="handleRegister(event)">
                        <div class="form-group">
                            <label class="form-label">用户名</label>
                            <input type="text" class="form-input" id="reg-username" required placeholder="请输入用户名">
                        </div>
                        <div class="form-group">
                            <label class="form-label">邮箱</label>
                            <input type="email" class="form-input" id="reg-email" required placeholder="请输入邮箱">
                        </div>
                        <div class="form-group">
                            <label class="form-label">手机号</label>
                            <input type="tel" class="form-input" id="reg-phone" placeholder="请输入手机号">
                        </div>
                        <div class="form-group">
                            <label class="form-label">密码</label>
                            <input type="password" class="form-input" id="reg-password" required placeholder="请输入密码">
                        </div>
                        <div class="form-group">
                            <label class="form-label">账号类型</label>
                            <select class="form-input form-select" id="reg-role">
                                <option value="customer">普通用户</option>
                                <option value="seller">商家</option>
                            </select>
                        </div>
                        <button type="submit" class="btn btn-primary btn-block">注册</button>
                    </form>
                    <p class="text-center mt-2 text-muted">
                        已有账号？<a href="/login" class="text-primary">立即登录</a>
                    </p>
                </div>
            </div>
        </div>
    `;
}

async function handleRegister(e) {
    e.preventDefault();
    
    const data = {
        username: $('#reg-username').value,
        email: $('#reg-email').value,
        phone: $('#reg-phone').value,
        password: $('#reg-password').value,
        role: $('#reg-role').value
    };

    const res = await api('/auth/register', {
        method: 'POST',
        body: data
    });

    if (res && res.code === 200) {
        App.token = res.data.token;
        App.user = res.data.user;
        localStorage.setItem('token', res.data.token);
        
        showToast('注册成功', 'success');
        updateHeader();
        navigate('/');
    } else {
        showToast(res?.message || '注册失败', 'error');
    }
}

// ========================================
// 商品详情页
// ========================================
async function renderProductDetail(id) {
    const content = $('#app-content');
    const res = await api('/product?id=' + id);
    
    if (!res || res.code !== 200) {
        content.innerHTML = '<div class="empty-state"><div class="empty-icon">😕</div><p class="empty-text">商品不存在</p></div>';
        return;
    }

    const product = res.data;
    const images = product.images ? product.images.split(',') : ['https://via.placeholder.com/400'];

    content.innerHTML = `
        <div class="product-detail">
            <div class="product-gallery">
                <div class="gallery-main">
                    <img src="${images[0]}" alt="${product.name}" id="main-image">
                </div>
                ${images.length > 1 ? `
                    <div class="gallery-thumbs">
                        ${images.map((img, i) => `
                            <div class="gallery-thumb ${i === 0 ? 'active' : ''}" onclick="changeImage('${img}', this)">
                                <img src="${img}" alt="">
                            </div>
                        `).join('')}
                    </div>
                ` : ''}
            </div>
            <div class="product-detail-info">
                <h1>${product.name}</h1>
                <div class="detail-price">
                    <span class="current">${formatPrice(product.price)}</span>
                    ${product.original_price > product.price ? 
                        `<span class="original">${formatPrice(product.original_price)}</span>` : ''}
                </div>
                <div class="detail-meta">
                    <span>⭐ ${(product.rating || 5).toFixed(1)} 分</span>
                    <span>📦 已售 ${product.sales || 0}</span>
                    <span>📊 库存 ${product.stock}</span>
                </div>
                <div class="mb-2">
                    <p class="text-muted">店铺：${product.seller_name || '官方自营'}</p>
                    <p class="text-muted">分类：${product.category_name || '未分类'}</p>
                </div>
                <div class="mb-3">
                    <label class="form-label">购买数量</label>
                    <div class="quantity-control">
                        <button class="quantity-btn" onclick="changeQuantity(-1)">−</button>
                        <input type="number" class="quantity-input" id="buy-quantity" value="1" min="1" max="${product.stock}">
                        <button class="quantity-btn" onclick="changeQuantity(1)">+</button>
                    </div>
                </div>
                <div class="detail-actions">
                    <button class="btn btn-outline btn-lg" onclick="addToCart(${product.id})">
                        🛒 加入购物车
                    </button>
                    <button class="btn btn-primary btn-lg" onclick="buyNow(${product.id})">
                        ⚡ 立即购买
                    </button>
                </div>
            </div>
        </div>

        <div class="card mt-3">
            <div class="card-header">
                <h3 class="card-title">商品详情</h3>
            </div>
            <div class="card-body">
                <p>${product.description || '暂无详情'}</p>
            </div>
        </div>

        <div class="card mt-3">
            <div class="card-header">
                <h3 class="card-title">商品评价</h3>
            </div>
            <div class="card-body" id="reviews-container">
                <div class="loading"><div class="spinner"></div></div>
            </div>
        </div>
    `;

    // 加载评价
    loadProductReviews(id);
}

function changeImage(src, thumb) {
    $('#main-image').src = src;
    $$('.gallery-thumb').forEach(el => el.classList.remove('active'));
    thumb.classList.add('active');
}

function changeQuantity(delta) {
    const input = $('#buy-quantity');
    const newVal = Math.max(1, Math.min(parseInt(input.max), parseInt(input.value) + delta));
    input.value = newVal;
}

async function addToCart(productId) {
    if (!App.user) {
        showToast('请先登录', 'warning');
        navigate('/login');
        return;
    }

    const quantity = parseInt($('#buy-quantity')?.value || 1);
    const res = await api('/cart/add', {
        method: 'POST',
        body: { product_id: productId, quantity }
    });

    if (res && res.code === 200) {
        showToast('已添加到购物车', 'success');
        updateCartCount();
    } else {
        showToast(res?.message || '添加失败', 'error');
    }
}

async function buyNow(productId) {
    if (!App.user) {
        showToast('请先登录', 'warning');
        navigate('/login');
        return;
    }

    await addToCart(productId);
    navigate('/cart');
}

async function loadProductReviews(productId) {
    const container = $('#reviews-container');
    const res = await api(`/reviews?product_id=${productId}&page=1&size=10`);

    if (res && res.code === 200 && res.data && res.data.length > 0) {
        container.innerHTML = res.data.map(review => `
            <div class="review-item">
                <div class="review-header">
                    <div class="review-avatar">${review.username?.charAt(0) || '?'}</div>
                    <div>
                        <div class="review-user">${review.username || '匿名用户'}</div>
                        <div class="review-rating">${'★'.repeat(review.rating)}${'☆'.repeat(5-review.rating)}</div>
                    </div>
                </div>
                <p class="review-content">${review.content || '用户未填写评价'}</p>
                ${review.reply ? `
                    <div class="review-reply">
                        <div class="review-reply-label">商家回复：</div>
                        <p>${review.reply}</p>
                    </div>
                ` : ''}
            </div>
        `).join('');
    } else {
        container.innerHTML = '<p class="text-muted text-center">暂无评价</p>';
    }
}

// ========================================
// 购物车页面
// ========================================
async function renderCart() {
    if (!App.user) {
        navigate('/login');
        return;
    }

    const content = $('#app-content');
    const res = await api('/cart');

    if (!res || res.code !== 200) {
        content.innerHTML = '<div class="empty-state"><div class="empty-icon">🛒</div><p class="empty-text">购物车加载失败</p></div>';
        return;
    }

    const { items, total_price, selected_count } = res.data;

    if (!items || items.length === 0) {
        content.innerHTML = `
            <div class="empty-state">
                <div class="empty-icon">🛒</div>
                <p class="empty-text">购物车是空的</p>
                <button class="btn btn-primary" onclick="navigate('/')">去购物</button>
            </div>
        `;
        return;
    }

    content.innerHTML = `
        <div class="card">
            <div class="card-header">
                <div class="d-flex align-center gap-2">
                    <input type="checkbox" class="cart-checkbox" id="select-all" onchange="selectAllCart(this.checked)" ${items.every(i => i.selected) ? 'checked' : ''}>
                    <label for="select-all">全选</label>
                </div>
                <button class="btn btn-sm btn-secondary" onclick="clearCart()">清空购物车</button>
            </div>
            <div class="card-body" style="padding: 0;">
                ${items.map(item => `
                    <div class="cart-item" data-id="${item.id}">
                        <input type="checkbox" class="cart-checkbox" ${item.selected ? 'checked' : ''} onchange="selectCartItem(${item.id}, this.checked)">
                        <div class="cart-image" onclick="navigate('/product/${item.product_id}')">
                            <img src="${getFirstImage(item.product.images)}" alt="">
                        </div>
                        <div class="cart-details">
                            <h4 class="cart-name" onclick="navigate('/product/${item.product_id}')">${item.product.name}</h4>
                            <p class="cart-price">${formatPrice(item.product.price)}</p>
                            <p class="text-muted" style="font-size: 0.85rem;">店铺：${item.product.seller_name || '官方自营'}</p>
                        </div>
                        <div class="cart-actions">
                            <div class="quantity-control">
                                <button class="quantity-btn" onclick="updateCartQuantity(${item.id}, ${item.quantity - 1})">−</button>
                                <input type="number" class="quantity-input" value="${item.quantity}" min="1" onchange="updateCartQuantity(${item.id}, this.value)">
                                <button class="quantity-btn" onclick="updateCartQuantity(${item.id}, ${item.quantity + 1})">+</button>
                            </div>
                            <button class="cart-delete" onclick="deleteCartItem(${item.id})">删除</button>
                        </div>
                    </div>
                `).join('')}
            </div>
            <div class="card-footer d-flex align-center justify-between">
                <div>
                    已选 <span class="text-primary">${selected_count}</span> 件商品
                </div>
                <div class="d-flex align-center gap-3">
                    <div>
                        合计：<span class="text-danger" style="font-size: 1.5rem; font-weight: bold;">${formatPrice(total_price)}</span>
                    </div>
                    <button class="btn btn-primary btn-lg" onclick="goToCheckout()" ${selected_count === 0 ? 'disabled' : ''}>
                        去结算
                    </button>
                </div>
            </div>
        </div>
    `;
}

async function selectCartItem(id, selected) {
    await api('/cart/select', {
        method: 'POST',
        body: { id, selected: selected ? 1 : 0 }
    });
    renderCart();
}

async function selectAllCart(selected) {
    await api('/cart/select-all', {
        method: 'POST',
        body: { selected: selected ? 1 : 0 }
    });
    renderCart();
}

async function updateCartQuantity(id, quantity) {
    if (quantity < 1) {
        deleteCartItem(id);
        return;
    }
    await api('/cart/update', {
        method: 'POST',
        body: { id, quantity: parseInt(quantity) }
    });
    renderCart();
}

async function deleteCartItem(id) {
    await api('/cart/delete?id=' + id, { method: 'DELETE' });
    showToast('已删除', 'success');
    renderCart();
    updateCartCount();
}

async function clearCart() {
    if (!confirm('确定要清空购物车吗？')) return;
    await api('/cart/clear', { method: 'DELETE' });
    showToast('已清空', 'success');
    renderCart();
    updateCartCount();
}

function goToCheckout() {
    navigate('/checkout');
}

// ========================================
// 结算页面
// ========================================
async function renderCheckout() {
    if (!App.user) {
        navigate('/login');
        return;
    }

    const content = $('#app-content');
    
    // 获取地址和购物车
    const [addressRes, cartRes] = await Promise.all([
        api('/addresses'),
        api('/cart')
    ]);

    const addresses = addressRes?.data || [];
    const cartData = cartRes?.data || {};
    const items = (cartData.items || []).filter(i => i.selected);

    if (items.length === 0) {
        content.innerHTML = `
            <div class="empty-state">
                <div class="empty-icon">📦</div>
                <p class="empty-text">请先选择要购买的商品</p>
                <button class="btn btn-primary" onclick="navigate('/cart')">返回购物车</button>
            </div>
        `;
        return;
    }

    const defaultAddress = addresses.find(a => a.is_default) || addresses[0];

    content.innerHTML = `
        <div style="display: grid; grid-template-columns: 1fr 350px; gap: 20px;">
            <div>
                <div class="checkout-section">
                    <h3 class="checkout-title">📍 收货地址</h3>
                    ${addresses.length === 0 ? `
                        <p class="text-muted">暂无收货地址</p>
                        <button class="btn btn-outline" onclick="showAddAddressModal()">+ 添加地址</button>
                    ` : `
                        <div id="address-list">
                            ${addresses.map(addr => `
                                <div class="address-card ${addr.id === defaultAddress?.id ? 'selected' : ''}" 
                                     onclick="selectAddress(${addr.id}, this)" data-id="${addr.id}">
                                    <div class="d-flex justify-between mb-1">
                                        <strong>${addr.receiver_name}</strong>
                                        <span>${addr.phone}</span>
                                    </div>
                                    <p class="text-muted">${addr.province}${addr.city}${addr.district}${addr.detail_address}</p>
                                </div>
                            `).join('')}
                        </div>
                        <button class="btn btn-sm btn-secondary mt-2" onclick="showAddAddressModal()">+ 添加新地址</button>
                    `}
                </div>

                <div class="checkout-section">
                    <h3 class="checkout-title">📦 商品清单</h3>
                    ${items.map(item => `
                        <div class="order-product">
                            <div class="order-product-image">
                                <img src="${getFirstImage(item.product.images)}" alt="">
                            </div>
                            <div class="order-product-info">
                                <h4 class="order-product-name">${item.product.name}</h4>
                                <p class="order-product-price">${formatPrice(item.product.price)} × ${item.quantity}</p>
                            </div>
                            <div class="text-right">
                                <strong class="text-danger">${formatPrice(item.product.price * item.quantity)}</strong>
                            </div>
                        </div>
                    `).join('')}
                </div>

                <div class="checkout-section">
                    <h3 class="checkout-title">📝 订单备注</h3>
                    <textarea class="form-input form-textarea" id="order-remark" placeholder="选填，可以告诉商家您的特殊需求"></textarea>
                </div>
            </div>

            <div>
                <div class="checkout-section" style="position: sticky; top: 90px;">
                    <h3 class="checkout-title">订单汇总</h3>
                    <div class="summary-row">
                        <span>商品金额</span>
                        <span>${formatPrice(cartData.total_price)}</span>
                    </div>
                    <div class="summary-row">
                        <span>运费</span>
                        <span>免运费</span>
                    </div>
                    <div class="summary-row summary-total">
                        <span>应付金额</span>
                        <span>${formatPrice(cartData.total_price)}</span>
                    </div>
                    <button class="btn btn-primary btn-block btn-lg mt-3" onclick="submitOrder()" ${addresses.length === 0 ? 'disabled' : ''}>
                        提交订单
                    </button>
                </div>
            </div>
        </div>
    `;

    window.selectedAddressId = defaultAddress?.id;
}

function selectAddress(id, element) {
    $$('.address-card').forEach(el => el.classList.remove('selected'));
    element.classList.add('selected');
    window.selectedAddressId = id;
}

function showAddAddressModal() {
    showModal('添加收货地址', `
        <form id="address-form">
            <div class="form-group">
                <label class="form-label">收货人</label>
                <input type="text" class="form-input" id="addr-name" required placeholder="请输入收货人姓名">
            </div>
            <div class="form-group">
                <label class="form-label">手机号</label>
                <input type="tel" class="form-input" id="addr-phone" required placeholder="请输入手机号">
            </div>
            <div class="form-group">
                <label class="form-label">省份</label>
                <input type="text" class="form-input" id="addr-province" placeholder="如：广东省">
            </div>
            <div class="form-group">
                <label class="form-label">城市</label>
                <input type="text" class="form-input" id="addr-city" placeholder="如：深圳市">
            </div>
            <div class="form-group">
                <label class="form-label">区县</label>
                <input type="text" class="form-input" id="addr-district" placeholder="如：南山区">
            </div>
            <div class="form-group">
                <label class="form-label">详细地址</label>
                <input type="text" class="form-input" id="addr-detail" required placeholder="请输入详细地址">
            </div>
            <div class="form-group">
                <label>
                    <input type="checkbox" id="addr-default"> 设为默认地址
                </label>
            </div>
        </form>
    `, `
        <button class="btn btn-secondary" onclick="closeModal()">取消</button>
        <button class="btn btn-primary" onclick="saveAddress()">保存</button>
    `);
}

async function saveAddress() {
    const data = {
        receiver_name: $('#addr-name').value,
        phone: $('#addr-phone').value,
        province: $('#addr-province').value,
        city: $('#addr-city').value,
        district: $('#addr-district').value,
        detail_address: $('#addr-detail').value,
        is_default: $('#addr-default').checked ? 1 : 0
    };

    const res = await api('/address/create', {
        method: 'POST',
        body: data
    });

    if (res && res.code === 200) {
        showToast('地址添加成功', 'success');
        closeModal();
        renderCheckout();
    } else {
        showToast(res?.message || '添加失败', 'error');
    }
}

async function submitOrder() {
    if (!window.selectedAddressId) {
        showToast('请选择收货地址', 'warning');
        return;
    }

    const res = await api('/order/create', {
        method: 'POST',
        body: {
            address_id: window.selectedAddressId,
            remark: $('#order-remark')?.value || ''
        }
    });

    if (res && res.code === 200) {
        showToast('订单创建成功', 'success');
        updateCartCount();
        
        // 显示支付弹窗
        const orderIds = res.data.order_ids;
        showPaymentModal(orderIds);
    } else {
        showToast(res?.message || '创建失败', 'error');
    }
}

function showPaymentModal(orderIds) {
    showModal('订单支付', `
        <div class="text-center">
            <p class="mb-2">订单已创建成功！</p>
            <p class="text-muted">请选择支付方式</p>
            <div class="d-flex gap-2 justify-center mt-3">
                <button class="btn btn-outline" onclick="payOrders(${JSON.stringify(orderIds)}, 'alipay')">支付宝</button>
                <button class="btn btn-outline" onclick="payOrders(${JSON.stringify(orderIds)}, 'wechat')">微信支付</button>
            </div>
        </div>
    `);
}

async function payOrders(orderIds, payType) {
    for (const id of orderIds) {
        await api('/order/pay', {
            method: 'POST',
            body: { order_id: id, pay_type: payType }
        });
    }
    
    showToast('支付成功', 'success');
    closeModal();
    navigate('/orders');
}
