describe('Home Page', () => {
  beforeEach(() => {
    document.body.innerHTML = `
      <main class="fade-in">
        <h1>Welcome to AlGoRhythm!</h1>
      </main>
    `;
  });

  it('verify the home page welcome message', () => {
    expect(document.querySelector('h1').textContent).toBe('Welcome to AlGoRhythm!');
  });
});
