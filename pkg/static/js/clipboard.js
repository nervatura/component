document.addEventListener("click", function (e) {
    const btn = e.target.closest("[data-copy]");
    if (!btn) return;

    navigator.clipboard.writeText(btn.dataset.copy);
});