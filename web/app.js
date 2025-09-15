(() => {
  const els = {
    booksCount: document.getElementById('books-count'),
    hadithCount: document.getElementById('hadith-count'),
    form: document.getElementById('search-form'),
    q: document.getElementById('q'),
    book: document.getElementById('book'),
    limit: document.getElementById('limit'),
    list: document.getElementById('list'),
    summary: document.getElementById('results-summary'),
    prev: document.getElementById('prev'),
    next: document.getElementById('next'),
  };

  let state = {
    results: [],          // current page items from server
    total: 0,             // total hits reported by server
    page: 1,              // 1-based
    pageSize: 10,
  };

  // Helpers
  const h = (tag, props = {}, children = []) => {
    const el = document.createElement(tag);
    for (const [k, v] of Object.entries(props)) {
      if (k === 'class') el.className = v;
      else if (k === 'text') el.textContent = v;
      else el.setAttribute(k, v);
    }
    for (const c of [].concat(children)) {
      if (typeof c === 'string') el.appendChild(document.createTextNode(c));
      else if (c) el.appendChild(c);
    }
    return el;
  };

  const debounce = (fn, ms) => {
    let t; return (...args) => { clearTimeout(t); t = setTimeout(() => fn(...args), ms); };
  };

  const fmtSummary = () => {
    const total = state.total;
    const start = total ? ((state.page - 1) * state.pageSize + 1) : 0;
    const end = Math.min(total, state.page * state.pageSize);
    const q = els.q.value.trim();
    const b = els.book.value;
    if (total === 0) {
      els.summary.textContent = b && !q ? `Browsing ${b}: no items` : 'No results';
      return;
    }
    if (!q && b) {
      els.summary.textContent = `Browsing ${b}: ${start}–${end} of ${total}`;
    } else {
      els.summary.textContent = `Showing ${start}–${end} of ${total}`;
    }
  };

  const render = () => {
    // Pagination controls
    const total = state.total;
    els.prev.disabled = state.page <= 1;
    const maxPage = Math.max(1, Math.ceil(total / state.pageSize));
    els.next.disabled = state.page >= maxPage;
    fmtSummary();

    // List
    els.list.innerHTML = '';
    for (const item of state.results) {
      const hadith = item.hadith || item.Hadith || {}; // API may use capital key for outer field
      const score = item.score ?? item.Score ?? 0;
      const book = hadith.book || hadith.Book || '';
      const number = hadith.number ?? hadith.Number ?? '';
      const id = hadith.id || hadith.Id || '';
      const arab = hadith.arab || hadith.Arab || '';

      const head = h('div', { class: 'head' }, [
        h('div', { class: 'book', text: book }),
        h('div', { class: 'no', text: `#${number}` }),
        h('div', { class: 'score', text: score ? `score: ${score}` : '' }),
      ]);
      const idLine = h('div', { class: 'id' }, [id]);
      const arLine = h('div', { class: 'ar arabic', lang: 'ar' }, [arab]);
      const li = h('li', { class: 'item' }, [head, idLine, arLine]);
      els.list.appendChild(li);
    }
  };

  const refetch = async () => {
    const q = els.q.value.trim();
    const page = state.page;
    const pageSize = state.pageSize;
    const selectedBook = els.book.value;
    // If neither query nor book is specified, nothing to show
    if (!q && !selectedBook) { state.results = []; state.total = 0; render(); return; }
    const url = new URL('/search', window.location.origin);
    url.searchParams.set('q', q);
    url.searchParams.set('page', String(page));
    url.searchParams.set('page_size', String(pageSize));
    if (selectedBook) url.searchParams.set('book', selectedBook);
    const res = await fetch(url.toString());
    if (!res.ok) {
      state.results = [];
      state.total = 0;
      render();
      return;
    }
    const total = parseInt(res.headers.get('X-Total-Count') || '0', 10);
    const data = await res.json();
    state.results = Array.isArray(data) ? data : [];
    state.total = Number.isFinite(total) ? total : state.results.length;
    render();
  };

  const search = async () => { state.page = 1; await refetch(); };

  const initCounts = async () => {
    try {
      const [booksRes, countRes] = await Promise.all([
        fetch('/books'),
        fetch('/count')
      ]);
      const books = booksRes.ok ? await booksRes.json() : [];
      const count = countRes.ok ? await countRes.json() : { count: 0 };
      els.booksCount.textContent = Array.isArray(books) ? books.length : '0';
      els.hadithCount.textContent = typeof count.count === 'number' ? count.count : '0';
      // Populate select
      if (Array.isArray(books)) {
        for (const b of books) {
          const opt = document.createElement('option');
          opt.value = String(b);
          opt.textContent = String(b);
          els.book.appendChild(opt);
        }
      }
    } catch (_) {
      // ignore
    }
  };

  // Events
  els.form.addEventListener('submit', (e) => { e.preventDefault(); search(); });
  els.book.addEventListener('change', () => { state.page = 1; refetch(); });
  els.limit.addEventListener('change', () => { state.pageSize = parseInt(els.limit.value, 10) || 10; state.page = 1; refetch(); });
  els.prev.addEventListener('click', () => { if (state.page > 1) { state.page--; refetch(); } });
  els.next.addEventListener('click', () => { state.page++; refetch(); });
  els.q.addEventListener('input', debounce(search, 250));

  // Boot
  initCounts();
})();
