describe('Home Page', () => {
  beforeEach(() => {
    document.body.innerHTML = `
      <main class="fade-in">
        <h1>Welcome to AlGoRhythm!</h1>
        <h2 class="fade-in" style="color:#1DB954; font-size:1.3rem; margin-top:18px;">Hello, TestUser! Glad to see you back.</h2>
      </main>
    `;
  });

  it('verify the home page welcome message', () => {
    expect(document.querySelector('h1').textContent).toBe('Welcome to AlGoRhythm!');
  });
});
