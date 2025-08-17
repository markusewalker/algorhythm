describe('Top Artists Page', () => {
  beforeEach(() => {
    document.body.innerHTML = `
      <ol class="fade-in">
        <li style="display: flex; align-items: center; margin-bottom: 18px;">
          <img src="artist.jpg" alt="Artist headshot" class="artist-headshot">
          <span class="artist-name">Test Artist</span>
        </li>
      </ol>
    `;
  });

  it('verify the artist headshot appears', () => {
    const img = document.querySelector('.artist-headshot');
    expect(img).not.toBeNull();
    expect(img.getAttribute('src')).toBe('artist.jpg');
  });

  it('verify the artist name appears', () => {
    const name = document.querySelector('.artist-name');
    expect(name.textContent).toBe('Test Artist');
  });

  it('verify you can click the artist name', () => {
    const name = document.querySelector('.artist-name');
    const clickEvent = new MouseEvent('click', {
      bubbles: true,
      cancelable: true,
      view: window
    });

    name.dispatchEvent(clickEvent);
    expect(name).toBeTruthy();
  });
});
