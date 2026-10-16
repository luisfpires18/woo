/**
 * Title banner. PLACEHOLDER treatment until the owner supplies the approved clean forge/cavern banner asset:
 * a CSS-only forge scene (stone, lava glow, ember light) behind the real HTML title. To swap in the final art,
 * set --banner-art (tokens.css) to the image URL and drop the placeholder title if the art carries the wordmark.
 * The banner looks the same in light and dark appearance; artwork is never darkened for dark mode.
 */
export function Banner() {
  return (
    <section className="banner" aria-labelledby="banner-title">
      <div className="banner__scene" aria-hidden="true">
        <span className="banner__glow banner__glow--left" />
        <span className="banner__glow banner__glow--right" />
        <span className="banner__lava" />
      </div>
      <div className="banner__content">
        <h1 id="banner-title" className="banner__title">
          <span className="banner__title-main">Weapons</span>{' '}
          <span className="banner__title-of">of</span>{' '}
          <span className="banner__title-main">Order</span>
        </h1>
      </div>
      <p className="banner__note">Placeholder banner · final artwork pending</p>
    </section>
  )
}
