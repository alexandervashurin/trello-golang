let currentBoardId = null;
let currentListId = null;
let isLoginMode = true;
let draggedCardId = null;

function getToken() {
    return localStorage.getItem('token');
}

function setToken(token) {
    localStorage.setItem('token', token);
}

function clearToken() {
    localStorage.removeItem('token');
}

function api(path, options = {}) {
    const token = getToken();
    const headers = { 'Content-Type': 'application/json' };
    if (token) {
        headers['Authorization'] = 'Bearer ' + token;
    }
    return fetch(path, { ...options, headers }).then(r => {
        if (r.status === 401) {
            clearToken();
            showAuthPage();
            return Promise.reject('unauthorized');
        }
        return r.json();
    });
}

function closeModal(e) {
    if (e && e.target !== e.currentTarget) return;
    document.querySelectorAll('.modal-overlay').forEach(m => m.classList.remove('active'));
}

function showAuthPage() {
    document.getElementById('authPage').style.display = 'flex';
    document.getElementById('app').style.display = 'none';
}

function showApp(username) {
    document.getElementById('authPage').style.display = 'none';
    document.getElementById('app').style.display = 'block';
    document.getElementById('userDisplay').textContent = username || '';
}

function toggleAuthMode() {
    isLoginMode = !isLoginMode;
    if (isLoginMode) {
        document.getElementById('authTitle').textContent = 'Вход';
        document.getElementById('authSubmitBtn').textContent = 'Войти';
        document.getElementById('authUsername').style.display = 'none';
        document.getElementById('authSwitchText').textContent = 'Нет аккаунта?';
        document.getElementById('authSwitchLink').textContent = 'Зарегистрироваться';
    } else {
        document.getElementById('authTitle').textContent = 'Регистрация';
        document.getElementById('authSubmitBtn').textContent = 'Зарегистрироваться';
        document.getElementById('authUsername').style.display = 'block';
        document.getElementById('authSwitchText').textContent = 'Уже есть аккаунт?';
        document.getElementById('authSwitchLink').textContent = 'Войти';
    }
    document.getElementById('authError').classList.remove('visible');
}

function showAuthError(msg) {
    const el = document.getElementById('authError');
    el.textContent = msg;
    el.classList.add('visible');
}

function handleAuth() {
    const email = document.getElementById('authEmail').value.trim();
    const password = document.getElementById('authPassword').value.trim();

    if (!email || !password) {
        showAuthError('Заполните все поля');
        return;
    }

    if (isLoginMode) {
        api('/api/login', {
            method: 'POST',
            body: JSON.stringify({ email, password }),
        }).then(data => {
            if (data.error) {
                showAuthError(data.error);
                return;
            }
            setToken(data.token);
            showApp(data.user.username);
            renderBoards();
        }).catch(() => {});
    } else {
        const username = document.getElementById('authUsername').value.trim();
        if (!username) {
            showAuthError('Заполните все поля');
            return;
        }
        api('/api/register', {
            method: 'POST',
            body: JSON.stringify({ username, email, password }),
        }).then(data => {
            if (data.error) {
                showAuthError(data.error);
                return;
            }
            setToken(data.token);
            showApp(data.user.username);
            renderBoards();
        }).catch(() => {});
    }
}

function logout() {
    clearToken();
    showAuthPage();
    document.getElementById('authEmail').value = '';
    document.getElementById('authPassword').value = '';
    document.getElementById('authUsername').value = '';
    document.getElementById('authError').classList.remove('visible');
}

function showAddBoardModal() {
    document.getElementById('addBoardModal').classList.add('active');
    document.getElementById('boardNameInput').value = '';
    document.getElementById('boardDescInput').value = '';
    document.getElementById('boardPublicInput').checked = false;
    setTimeout(() => document.getElementById('boardNameInput').focus(), 100);
}

function showAddListModal() {
    document.getElementById('addListModal').classList.add('active');
    document.getElementById('listNameInput').value = '';
    setTimeout(() => document.getElementById('listNameInput').focus(), 100);
}

function showAddCardModal(listId) {
    currentListId = listId;
    document.getElementById('addCardModal').classList.add('active');
    document.getElementById('cardTitleInput').value = '';
    document.getElementById('cardDescInput').value = '';
    setTimeout(() => document.getElementById('cardTitleInput').focus(), 100);
}

function createBoard() {
    const name = document.getElementById('boardNameInput').value.trim();
    if (!name) return;
    api('/api/boards', {
        method: 'POST',
        body: JSON.stringify({
            name,
            description: document.getElementById('boardDescInput').value.trim(),
            is_public: document.getElementById('boardPublicInput').checked,
        }),
    }).then(() => {
        closeModal();
        renderBoards();
    });
}

function createList() {
    const name = document.getElementById('listNameInput').value.trim();
    if (!name || !currentBoardId) return;
    api('/api/lists', {
        method: 'POST',
        body: JSON.stringify({ name, board_id: currentBoardId, position: 0 }),
    }).then(() => {
        closeModal();
        renderBoard(currentBoardId);
    });
}

function createCard() {
    const title = document.getElementById('cardTitleInput').value.trim();
    if (!title || !currentListId) return;
    api('/api/cards', {
        method: 'POST',
        body: JSON.stringify({ title, description: document.getElementById('cardDescInput').value.trim(), list_id: currentListId, position: 0 }),
    }).then(() => {
        closeModal();
        renderBoard(currentBoardId);
    });
}

function deleteBoard(id, e) {
    e.stopPropagation();
    api(`/api/board?id=${id}`, { method: 'DELETE' }).then(() => renderBoards());
}

function deleteList(id, e) {
    e.stopPropagation();
    api(`/api/list?id=${id}`, { method: 'DELETE' }).then(() => renderBoard(currentBoardId));
}

function deleteCard(id, e) {
    e.stopPropagation();
    api(`/api/card?id=${id}`, { method: 'DELETE' }).then(() => renderBoard(currentBoardId));
}

function showBoards() {
    document.getElementById('boardsView').style.display = 'block';
    document.getElementById('boardView').style.display = 'none';
    currentBoardId = null;
    renderBoards();
}

function openBoard(id) {
    currentBoardId = id;
    document.getElementById('boardsView').style.display = 'none';
    document.getElementById('boardView').style.display = 'flex';
    renderBoard(id);
}

function renderBoards() {
    api('/api/boards').then(boards => {
        const grid = document.getElementById('boardsGrid');
        if (!boards || boards.length === 0) {
            grid.innerHTML = '<div class="empty-state"><h3>Нет досок</h3><p>Создайте первую доску, чтобы начать работу</p></div>';
        } else {
            grid.innerHTML = boards.map(b => `
                <div class="board-card" onclick="openBoard('${b.id}')">
                    <div>
                        <h3>${esc(b.name)} ${b.is_public ? '<span class="badge badge-public">публичная</span>' : '<span class="badge badge-private">частная</span>'}</h3>
                        ${b.description ? `<p>${esc(b.description)}</p>` : ''}
                    </div>
                    <div class="board-card-actions">
                        <button class="btn btn-danger" onclick="deleteBoard('${b.id}', event)">Удалить</button>
                    </div>
                </div>
            `).join('');
        }
    });

    api('/api/boards/public').then(boards => {
        const grid = document.getElementById('publicBoardsGrid');
        if (!boards || boards.length === 0) {
            grid.innerHTML = '<div class="empty-state"><p style="padding:20px 0;">Нет публичных досок</p></div>';
            return;
        }
        grid.innerHTML = boards.map(b => `
            <div class="board-card" onclick="openBoard('${b.id}')">
                <div>
                    <h3>${esc(b.name)} <span class="badge badge-public">публичная</span></h3>
                    ${b.description ? `<p>${esc(b.description)}</p>` : ''}
                </div>
            </div>
        `).join('');
    });
}

function renderBoard(boardId) {
    Promise.all([
        api(`/api/board?id=${boardId}`),
        api(`/api/lists?board_id=${boardId}`),
    ]).then(([board, lists]) => {
        document.getElementById('boardTitle').textContent = board.name;
        loadCards(lists);
    });
}

function loadCards(lists) {
    if (!lists || lists.length === 0) {
        document.getElementById('boardColumns').innerHTML = '<div class="empty-state"><h3>Нет списков</h3><p>Добавьте первый список в эту доску</p></div>';
        return;
    }

    const sortedLists = [...lists].sort((a, b) => a.position - b.position);
    const cardPromises = sortedLists.map(list =>
        api(`/api/cards?list_id=${list.id}`).then(cards => ({ list, cards }))
    );

    Promise.all(cardPromises).then(results => {
        const columns = document.getElementById('boardColumns');
        columns.innerHTML = results.map(({ list, cards }) => `
            <div class="board-column">
                <div class="column-header">
                    <span class="list-title">${esc(list.name)}</span>
                    <button class="btn btn-danger" onclick="deleteList('${list.id}', event)" style="padding:2px 8px;font-size:12px;">&times;</button>
                </div>
                <div class="column-cards" data-list-id="${list.id}"
                    ondrop="onCardDrop(event)"
                    ondragover="onDragOver(event)"
                    ondragleave="onDragLeave(event)">
                    ${(!cards || cards.length === 0)
                        ? '<div style="padding:12px;text-align:center;color:#a0a0b0;font-size:13px;">Нет карточек</div>'
                        : cards.sort((a, b) => a.position - b.position).map(c => `
                            <div class="card-item" draggable="true"
                                ondragstart="onCardDragStart(event, '${c.id}')"
                                ondragend="onCardDragEnd(event)">
                                <h4>${esc(c.title)}</h4>
                                ${c.description ? `<p>${esc(c.description)}</p>` : ''}
                                <div class="card-actions">
                                    <button class="btn btn-danger" onclick="deleteCard('${c.id}', event)" style="padding:2px 8px;font-size:12px;">Удалить</button>
                                </div>
                            </div>
                        `).join('')
                    }
                </div>
                <button class="add-card-btn" onclick="showAddCardModal('${list.id}')">+ Добавить карточку</button>
            </div>
        `).join('');
    });
}

function onCardDragStart(e, cardId) {
    draggedCardId = cardId;
    e.dataTransfer.effectAllowed = 'move';
    e.dataTransfer.setData('text/plain', cardId);
    e.currentTarget.classList.add('dragging');
}

function onCardDragEnd(e) {
    e.currentTarget.classList.remove('dragging');
    document.querySelectorAll('.column-cards').forEach(el => el.classList.remove('drag-over'));
    draggedCardId = null;
}

function onDragOver(e) {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
    e.currentTarget.classList.add('drag-over');
}

function onDragLeave(e) {
    e.currentTarget.classList.remove('drag-over');
}

function onCardDrop(e) {
    e.preventDefault();
    const target = e.currentTarget;
    target.classList.remove('drag-over');

    const cardId = e.dataTransfer.getData('text/plain') || draggedCardId;
    if (!cardId) return;

    const targetListId = target.getAttribute('data-list-id');
    if (!targetListId) return;

    const cardsInColumn = target.querySelectorAll('.card-item:not(.dragging)');
    const position = cardsInColumn.length;

    api('/api/card', {
        method: 'PATCH',
        body: JSON.stringify({ id: cardId, list_id: targetListId, position }),
    }).then(() => {
        renderBoard(currentBoardId);
    });
}

function esc(str) {
    const div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
}

const token = getToken();
if (token) {
    api('/api/me').then(user => {
        if (user && user.username) {
            showApp(user.username);
            renderBoards();
        } else {
            showAuthPage();
        }
    }).catch(() => showAuthPage());
} else {
    showAuthPage();
}
