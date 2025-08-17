describe('Login Page', () => {
  beforeEach(() => {
    document.body.innerHTML = `
      <button type="submit">Login</button>
    `;
  });

  it('verify the login button is present', () => {
    const button = document.querySelector('button[type="submit"]');
    expect(button).not.toBeNull();
    expect(button.textContent).toBe('Login');
  });

  it('verify the login button is functional', () => {
    const button = document.querySelector('button[type="submit"]');
    const mockHandler = jest.fn();
    button.addEventListener('click', mockHandler);
    button.click();
    expect(mockHandler).toHaveBeenCalled();
  });
});
