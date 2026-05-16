(() => {
  'use strict';

  const DEFAULT_API_BASE = 'http://127.0.0.1:8080/api/v1';
  const STORAGE_KEY = 'student-system-api-base';

  const state = {
    apiBase: localStorage.getItem(STORAGE_KEY) || DEFAULT_API_BASE,
    page: 1,
    pageSize: 20,
    total: 0,
    filters: {},
    deleteTargetId: null,
    deleteTargetName: '',
  };

  const els = {};
  let studentFormModal;
  let studentDetailModal;
  let deleteConfirmModal;
  let appToast;

  document.addEventListener('DOMContentLoaded', () => {
    cacheElements();
    initBootstrapComponents();
    bindEvents();
    restoreUiState();
    checkHealth(false);
    loadStudents();
  });

  function cacheElements() {
    Object.assign(els, {
      apiBaseInput: document.getElementById('apiBaseInput'),
      saveApiBaseBtn: document.getElementById('saveApiBaseBtn'),
      healthCheckBtn: document.getElementById('healthCheckBtn'),
      reloadPageBtn: document.getElementById('reloadPageBtn'),
      serviceStatusDot: document.getElementById('serviceStatusDot'),
      serviceStatusText: document.getElementById('serviceStatusText'),
      openCreateModalBtn: document.getElementById('openCreateModalBtn'),
      filterForm: document.getElementById('filterForm'),
      resetFilterBtn: document.getElementById('resetFilterBtn'),
      filterStudentNo: document.getElementById('filterStudentNo'),
      filterName: document.getElementById('filterName'),
      filterClassName: document.getElementById('filterClassName'),
      filterMajorName: document.getElementById('filterMajorName'),
      filterGradeYear: document.getElementById('filterGradeYear'),
      filterClassId: document.getElementById('filterClassId'),
      filterMajorId: document.getElementById('filterMajorId'),
      filterStatus: document.getElementById('filterStatus'),
      pageSizeSelect: document.getElementById('pageSizeSelect'),
      studentTableBody: document.getElementById('studentTableBody'),
      tableSummary: document.getElementById('tableSummary'),
      currentQueryBadge: document.getElementById('currentQueryBadge'),
      paginationText: document.getElementById('paginationText'),
      paginationList: document.getElementById('paginationList'),
      quickStudentNoInput: document.getElementById('quickStudentNoInput'),
      quickStudentNoBtn: document.getElementById('quickStudentNoBtn'),
      requestMeta: document.getElementById('requestMeta'),
      responsePreview: document.getElementById('responsePreview'),
      clearLogBtn: document.getElementById('clearLogBtn'),
      studentForm: document.getElementById('studentForm'),
      studentFormModalTitle: document.getElementById('studentFormModalTitle'),
      editingStudentId: document.getElementById('editingStudentId'),
      formStudentNo: document.getElementById('formStudentNo'),
      formName: document.getElementById('formName'),
      formClassId: document.getElementById('formClassId'),
      formStatusGroup: document.getElementById('formStatusGroup'),
      formStatus: document.getElementById('formStatus'),
      formPhone: document.getElementById('formPhone'),
      formEmail: document.getElementById('formEmail'),
      formAddress: document.getElementById('formAddress'),
      submitStudentFormBtn: document.getElementById('submitStudentFormBtn'),
      studentDetailContent: document.getElementById('studentDetailContent'),
      deleteTargetText: document.getElementById('deleteTargetText'),
      confirmDeleteBtn: document.getElementById('confirmDeleteBtn'),
      appToast: document.getElementById('appToast'),
      toastTitle: document.getElementById('toastTitle'),
      toastBody: document.getElementById('toastBody'),
      toastTime: document.getElementById('toastTime'),
    });
  }

  function initBootstrapComponents() {
    studentFormModal = new bootstrap.Modal(document.getElementById('studentFormModal'));
    studentDetailModal = new bootstrap.Modal(document.getElementById('studentDetailModal'));
    deleteConfirmModal = new bootstrap.Modal(document.getElementById('deleteConfirmModal'));
    appToast = new bootstrap.Toast(els.appToast, { delay: 3200 });
  }

  function bindEvents() {
    els.saveApiBaseBtn.addEventListener('click', saveApiBase);
    els.healthCheckBtn.addEventListener('click', () => checkHealth(true));
    els.reloadPageBtn.addEventListener('click', () => loadStudents());
    els.openCreateModalBtn.addEventListener('click', openCreateModal);
    els.filterForm.addEventListener('submit', handleFilterSubmit);
    els.resetFilterBtn.addEventListener('click', resetFilters);
    els.pageSizeSelect.addEventListener('change', handlePageSizeChange);
    els.quickStudentNoBtn.addEventListener('click', quickFindByStudentNo);
    els.quickStudentNoInput.addEventListener('keydown', (event) => {
      if (event.key === 'Enter') {
        quickFindByStudentNo();
      }
    });
    els.clearLogBtn.addEventListener('click', clearRequestLog);
    els.studentForm.addEventListener('submit', submitStudentForm);
    els.confirmDeleteBtn.addEventListener('click', confirmDeleteStudent);
    els.studentTableBody.addEventListener('click', handleTableAction);
  }

  function restoreUiState() {
    els.apiBaseInput.value = state.apiBase;
    els.pageSizeSelect.value = String(state.pageSize);
  }

  function normalizeApiBase(value) {
    return value.trim().replace(/\/+$/, '');
  }

  function saveApiBase() {
    const nextBase = normalizeApiBase(els.apiBaseInput.value);
    if (!nextBase) {
      showToast('地址错误', 'API 基础地址不能为空。', 'warning');
      return;
    }
    state.apiBase = nextBase;
    localStorage.setItem(STORAGE_KEY, nextBase);
    showToast('已保存', '后端 API 基础地址已更新。', 'success');
    checkHealth(true);
    loadStudents();
  }

  async function checkHealth(showResultToast) {
    try {
      const data = await request('/health', { method: 'GET' });
      setServiceStatus(true, data?.status === 'ok' ? '服务正常' : '服务已响应');
      if (showResultToast) {
        showToast('检测成功', '后端服务可正常访问。', 'success');
      }
    } catch (error) {
      setServiceStatus(false, '服务不可达');
      if (showResultToast) {
        showToast('检测失败', error.message, 'danger');
      }
    }
  }

  function setServiceStatus(ok, text) {
    els.serviceStatusDot.classList.remove('ok', 'fail');
    els.serviceStatusDot.classList.add(ok ? 'ok' : 'fail');
    els.serviceStatusText.textContent = text;
  }

  function handleFilterSubmit(event) {
    event.preventDefault();
    state.page = 1;
    state.filters = collectFilters();
    loadStudents();
  }

  function handlePageSizeChange() {
    state.pageSize = Number(els.pageSizeSelect.value) || 20;
    state.page = 1;
    loadStudents();
  }

  function collectFilters() {
    return compactObject({
      student_no: els.filterStudentNo.value.trim(),
      name: els.filterName.value.trim(),
      class_name: els.filterClassName.value.trim(),
      major_name: els.filterMajorName.value.trim(),
      grade_year: els.filterGradeYear.value.trim(),
      class_id: els.filterClassId.value.trim(),
      major_id: els.filterMajorId.value.trim(),
      status: els.filterStatus.value.trim(),
    });
  }

  function resetFilters() {
    els.filterForm.reset();
    els.pageSizeSelect.value = '20';
    state.page = 1;
    state.pageSize = 20;
    state.filters = {};
    loadStudents();
  }

  async function loadStudents() {
    setTableLoading();
    const params = new URLSearchParams({
      ...state.filters,
      page: String(state.page),
      page_size: String(state.pageSize),
    });

    try {
      const pageData = await request(`/students?${params.toString()}`, { method: 'GET' });
      state.total = Number(pageData.total || 0);
      state.page = Number(pageData.page || state.page);
      state.pageSize = Number(pageData.page_size || state.pageSize);
      renderStudentTable(Array.isArray(pageData.list) ? pageData.list : []);
      renderPagination();
      updateListSummary();
      updateQueryBadge();
    } catch (error) {
      renderTableError(error.message);
      renderPagination(true);
      updateListSummary('加载失败');
      showToast('列表加载失败', error.message, 'danger');
    }
  }

  function setTableLoading() {
    els.studentTableBody.innerHTML = `
      <tr>
        <td colspan="9" class="text-center py-5">
          <div class="spinner-border text-primary" role="status" aria-hidden="true"></div>
          <div class="text-secondary mt-3">正在加载学生数据...</div>
        </td>
      </tr>
    `;
  }

  function renderTableError(message) {
    els.studentTableBody.innerHTML = `
      <tr>
        <td colspan="9" class="text-center py-5 text-danger">
          ${escapeHtml(message)}
        </td>
      </tr>
    `;
  }

  function renderStudentTable(list) {
    if (!list.length) {
      els.studentTableBody.innerHTML = `
        <tr>
          <td colspan="9">
            <div class="empty-state">没有符合条件的学生数据</div>
          </td>
        </tr>
      `;
      return;
    }

    els.studentTableBody.innerHTML = list.map((student) => {
      const statusBadge = Number(student.status) === 1
        ? '<span class="badge text-bg-success">在读</span>'
        : '<span class="badge text-bg-secondary">停用</span>';
      const contact = [student.phone, student.email].filter(Boolean);
      const contactHtml = contact.length
        ? contact.map((item) => `<span>${escapeHtml(item)}</span>`).join('')
        : '<span class="text-secondary">未填写</span>';

      return `
        <tr>
          <td>${escapeHtml(student.id)}</td>
          <td><strong>${escapeHtml(student.student_no)}</strong></td>
          <td>${escapeHtml(student.name)}</td>
          <td>
            <div class="fw-semibold">${escapeHtml(student.major_name || '-')}</div>
            <div class="text-secondary small">${escapeHtml(student.class_name || '-')} · ID ${escapeHtml(student.class_id)}</div>
          </td>
          <td>${escapeHtml(student.grade_year || '-')}</td>
          <td><div class="contact-stack">${contactHtml}</div></td>
          <td>${statusBadge}</td>
          <td>${escapeHtml(formatDateTime(student.updated_at))}</td>
          <td class="text-end">
            <div class="action-group">
              <button type="button" class="btn btn-sm btn-outline-primary" data-action="detail" data-id="${escapeHtml(student.id)}">详情</button>
              <button type="button" class="btn btn-sm btn-outline-secondary" data-action="edit" data-id="${escapeHtml(student.id)}">编辑</button>
              <button type="button" class="btn btn-sm btn-outline-danger" data-action="delete" data-id="${escapeHtml(student.id)}" data-name="${escapeHtml(student.name)}">停用</button>
            </div>
          </td>
        </tr>
      `;
    }).join('');
  }

  function updateListSummary(forceText) {
    if (forceText) {
      els.tableSummary.textContent = forceText;
      return;
    }
    const start = state.total === 0 ? 0 : (state.page - 1) * state.pageSize + 1;
    const end = Math.min(state.total, state.page * state.pageSize);
    els.tableSummary.textContent = `共 ${state.total} 条，当前显示 ${start}-${end} 条`;
  }

  function updateQueryBadge() {
    const hasFilters = Object.keys(state.filters).length > 0;
    els.currentQueryBadge.textContent = hasFilters ? '已应用组合检索' : '默认在读学生';
  }

  function renderPagination(disabled = false) {
    if (disabled) {
      els.paginationText.textContent = '第 1 页';
      els.paginationList.innerHTML = '';
      return;
    }

    const totalPages = Math.max(1, Math.ceil(state.total / state.pageSize));
    const currentPage = Math.min(Math.max(state.page, 1), totalPages);
    state.page = currentPage;
    els.paginationText.textContent = `第 ${currentPage} / ${totalPages} 页`;

    const pages = buildPageRange(currentPage, totalPages);
    const items = [];
    items.push(createPageItem('上一页', currentPage - 1, currentPage <= 1));
    pages.forEach((page) => {
      if (page === '...') {
        items.push('<li class="page-item disabled"><span class="page-link">...</span></li>');
      } else {
        items.push(createPageItem(String(page), page, false, page === currentPage));
      }
    });
    items.push(createPageItem('下一页', currentPage + 1, currentPage >= totalPages));
    els.paginationList.innerHTML = items.join('');

    els.paginationList.querySelectorAll('[data-page]').forEach((link) => {
      link.addEventListener('click', () => {
        const nextPage = Number(link.dataset.page);
        if (!Number.isFinite(nextPage) || nextPage < 1 || nextPage > totalPages || nextPage === state.page) {
          return;
        }
        state.page = nextPage;
        loadStudents();
      });
    });
  }

  function buildPageRange(current, total) {
    if (total <= 7) {
      return Array.from({ length: total }, (_, index) => index + 1);
    }

    const pages = [1];
    const left = Math.max(2, current - 1);
    const right = Math.min(total - 1, current + 1);

    if (left > 2) pages.push('...');
    for (let page = left; page <= right; page += 1) pages.push(page);
    if (right < total - 1) pages.push('...');
    pages.push(total);
    return pages;
  }

  function createPageItem(label, page, disabled = false, active = false) {
    const classes = ['page-item'];
    if (disabled) classes.push('disabled');
    if (active) classes.push('active');
    const dataAttr = disabled || active ? '' : ` data-page="${page}"`;
    return `<li class="${classes.join(' ')}"><span class="page-link"${dataAttr}>${escapeHtml(label)}</span></li>`;
  }

  async function quickFindByStudentNo() {
    const studentNo = els.quickStudentNoInput.value.trim();
    if (!studentNo) {
      showToast('请输入学号', '学号快速定位需要完整学号。', 'warning');
      return;
    }

    try {
      const student = await request(`/students/by-no/${encodeURIComponent(studentNo)}`, { method: 'GET' });
      renderStudentDetail(student);
      studentDetailModal.show();
    } catch (error) {
      showToast('查询失败', error.message, 'danger');
    }
  }

  async function handleTableAction(event) {
    const button = event.target.closest('button[data-action]');
    if (!button) return;

    const id = Number(button.dataset.id);
    if (!Number.isFinite(id) || id <= 0) return;

    switch (button.dataset.action) {
      case 'detail':
        await openDetailModal(id);
        break;
      case 'edit':
        await openEditModal(id);
        break;
      case 'delete':
        openDeleteModal(id, button.dataset.name || '该学生');
        break;
      default:
        break;
    }
  }

  function openCreateModal() {
    els.studentForm.reset();
    els.editingStudentId.value = '';
    els.studentFormModalTitle.textContent = '新增学生';
    els.submitStudentFormBtn.textContent = '创建学生';
    els.formStudentNo.disabled = false;
    els.formStudentNo.required = true;
    els.formStatusGroup.classList.add('d-none');
    els.formStatus.value = '1';
    studentFormModal.show();
  }

  async function openEditModal(id) {
    try {
      const student = await request(`/students/${id}`, { method: 'GET' });
      els.studentForm.reset();
      els.editingStudentId.value = String(student.id);
      els.studentFormModalTitle.textContent = `编辑学生：${student.name}`;
      els.submitStudentFormBtn.textContent = '保存修改';
      els.formStudentNo.value = student.student_no || '';
      els.formStudentNo.disabled = true;
      els.formStudentNo.required = false;
      els.formName.value = student.name || '';
      els.formClassId.value = student.class_id || '';
      els.formStatusGroup.classList.remove('d-none');
      els.formStatus.value = String(student.status ?? 1);
      els.formPhone.value = student.phone || '';
      els.formEmail.value = student.email || '';
      els.formAddress.value = student.address || '';
      studentFormModal.show();
    } catch (error) {
      showToast('编辑失败', error.message, 'danger');
    }
  }

  async function openDetailModal(id) {
    try {
      const student = await request(`/students/${id}`, { method: 'GET' });
      renderStudentDetail(student);
      studentDetailModal.show();
    } catch (error) {
      showToast('详情加载失败', error.message, 'danger');
    }
  }

  function renderStudentDetail(student) {
    const detailItems = [
      ['学生 ID', student.id],
      ['学号', student.student_no],
      ['姓名', student.name],
      ['状态', Number(student.status) === 1 ? '在读 / 有效' : '停用 / 已删除'],
      ['专业', `${student.major_name || '-'}（${student.major_code || '-'} / ID ${student.major_id || '-'}）`],
      ['班级', `${student.class_name || '-'}（${student.class_code || '-'} / ID ${student.class_id || '-'}）`],
      ['年级', student.grade_year || '-'],
      ['联系电话', student.phone || '-'],
      ['邮箱', student.email || '-'],
      ['创建时间', formatDateTime(student.created_at)],
      ['更新时间', formatDateTime(student.updated_at)],
    ];

    const regularHtml = detailItems.map(([label, value]) => detailItem(label, value)).join('');
    const addressHtml = detailItem('地址', student.address || '-', true);
    els.studentDetailContent.innerHTML = regularHtml + addressHtml;
  }

  function detailItem(label, value, full = false) {
    return `
      <div class="detail-item${full ? ' full' : ''}">
        <div class="detail-label">${escapeHtml(label)}</div>
        <p class="detail-value">${escapeHtml(value)}</p>
      </div>
    `;
  }

  async function submitStudentForm(event) {
    event.preventDefault();
    const editingId = Number(els.editingStudentId.value || 0);
    const isEdit = Number.isFinite(editingId) && editingId > 0;

    const commonPayload = {
      name: els.formName.value.trim(),
      class_id: Number(els.formClassId.value),
      phone: els.formPhone.value.trim(),
      email: els.formEmail.value.trim(),
      address: els.formAddress.value.trim(),
    };

    if (!commonPayload.name) {
      showToast('表单未完成', '姓名不能为空。', 'warning');
      return;
    }
    if (!Number.isFinite(commonPayload.class_id) || commonPayload.class_id <= 0) {
      showToast('表单未完成', '班级 ID 必须为正整数。', 'warning');
      return;
    }

    try {
      setSubmitButtonLoading(true);
      if (isEdit) {
        const payload = {
          ...commonPayload,
          status: Number(els.formStatus.value),
        };
        await request(`/students/${editingId}`, {
          method: 'PUT',
          body: JSON.stringify(payload),
        });
        studentFormModal.hide();
        showToast('修改成功', '学生档案已更新。', 'success');
      } else {
        const studentNo = els.formStudentNo.value.trim();
        if (!studentNo) {
          showToast('表单未完成', '学号不能为空。', 'warning');
          return;
        }
        const payload = {
          student_no: studentNo,
          ...commonPayload,
        };
        await request('/students', {
          method: 'POST',
          body: JSON.stringify(payload),
        });
        studentFormModal.hide();
        showToast('创建成功', '学生档案已新增。', 'success');
      }
      loadStudents();
    } catch (error) {
      showToast(isEdit ? '修改失败' : '创建失败', error.message, 'danger');
    } finally {
      setSubmitButtonLoading(false);
    }
  }

  function setSubmitButtonLoading(loading) {
    els.submitStudentFormBtn.disabled = loading;
    els.submitStudentFormBtn.textContent = loading
      ? '处理中...'
      : (els.editingStudentId.value ? '保存修改' : '创建学生');
  }

  function openDeleteModal(id, name) {
    state.deleteTargetId = id;
    state.deleteTargetName = name;
    els.deleteTargetText.textContent = `目标学生：${name}（ID ${id}）`;
    deleteConfirmModal.show();
  }

  async function confirmDeleteStudent() {
    if (!state.deleteTargetId) return;
    try {
      els.confirmDeleteBtn.disabled = true;
      els.confirmDeleteBtn.textContent = '处理中...';
      await request(`/students/${state.deleteTargetId}`, { method: 'DELETE' });
      deleteConfirmModal.hide();
      showToast('停用成功', `${state.deleteTargetName} 已被设置为停用状态。`, 'success');
      loadStudents();
    } catch (error) {
      showToast('停用失败', error.message, 'danger');
    } finally {
      els.confirmDeleteBtn.disabled = false;
      els.confirmDeleteBtn.textContent = '确认停用';
      state.deleteTargetId = null;
      state.deleteTargetName = '';
    }
  }

  async function request(path, options = {}) {
    const url = `${state.apiBase}${path}`;
    const method = options.method || 'GET';
    const headers = new Headers(options.headers || {});
    if (options.body && !headers.has('Content-Type')) {
      headers.set('Content-Type', 'application/json');
    }

    const startedAt = performance.now();
    try {
      const response = await fetch(url, {
        ...options,
        method,
        headers,
      });
      const text = await response.text();
      const payload = safeJsonParse(text);
      const elapsed = Math.round(performance.now() - startedAt);
      updateRequestLog(method, url, response.status, payload ?? text, elapsed);

      if (!response.ok) {
        const message = payload?.message || `HTTP ${response.status}`;
        throw new Error(message);
      }
      if (!payload || typeof payload.code !== 'number') {
        throw new Error('接口响应格式不符合约定。');
      }
      if (payload.code !== 0) {
        throw new Error(payload.message || '接口返回业务错误。');
      }
      return payload.data;
    } catch (error) {
      if (!(error instanceof Error)) {
        throw new Error('未知网络错误。');
      }
      if (error.name === 'TypeError') {
        updateRequestLog(method, url, 'NETWORK_ERROR', { message: error.message }, Math.round(performance.now() - startedAt));
        throw new Error('网络请求失败，请检查后端服务与跨域配置。');
      }
      throw error;
    }
  }

  function updateRequestLog(method, url, status, payload, elapsed) {
    els.requestMeta.textContent = `${method} ${url} · ${status} · ${elapsed}ms`;
    els.responsePreview.textContent = typeof payload === 'string'
      ? payload
      : JSON.stringify(payload, null, 2);
  }

  function clearRequestLog() {
    els.requestMeta.textContent = '暂无请求';
    els.responsePreview.textContent = '等待接口调用...';
  }

  function showToast(title, body, variant = 'primary') {
    els.toastTitle.textContent = title;
    els.toastBody.textContent = body;
    els.toastTime.textContent = '刚刚';
    els.appToast.className = `toast border-0 shadow text-bg-${variant === 'success' ? 'success' : variant === 'danger' ? 'danger' : variant === 'warning' ? 'warning' : 'primary'}`;
    appToast.show();
  }

  function compactObject(input) {
    return Object.fromEntries(
      Object.entries(input).filter(([, value]) => value !== '' && value !== null && value !== undefined)
    );
  }

  function safeJsonParse(text) {
    if (!text) return null;
    try {
      return JSON.parse(text);
    } catch (_error) {
      return null;
    }
  }

  function formatDateTime(value) {
    if (!value) return '-';
    return String(value).replace('T', ' ');
  }

  function escapeHtml(value) {
    return String(value ?? '')
      .replaceAll('&', '&amp;')
      .replaceAll('<', '&lt;')
      .replaceAll('>', '&gt;')
      .replaceAll('"', '&quot;')
      .replaceAll("'", '&#039;');
  }
})();
