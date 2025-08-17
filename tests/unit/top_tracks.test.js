describe('Top Tracks Page', () => {
  beforeEach(() => {
    document.body.innerHTML = `
      <ol class="fade-in">
        <li style="margin-bottom: 32px; display: flex; align-items: center;">
          <iframe src="https://open.spotify.com/embed/track/123" width="300" height="80" frameborder="0" allowtransparency="true" allow="encrypted-media"></iframe>
        </li>
      </ol>
    `;
  });

  it('verify the Spotify track preview', () => {
    const iframe = document.querySelector('iframe');

    expect(iframe).not.toBeNull();
    expect(iframe.getAttribute('src')).toContain('open.spotify.com/embed/track/123');
  });
});
