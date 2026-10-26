const PageIntro = ({ eyebrow, title, children, meta }) => (
  <header className="page-intro">
    <div className="page-intro__leader" aria-hidden="true">24 FPS</div>
    <div className="page-intro__content">
      <p className="eyebrow">{eyebrow}</p>
      <h1>{title}</h1>
      {children && <p className="page-intro__copy">{children}</p>}
      {meta && <p className="page-intro__meta">{meta}</p>}
    </div>
  </header>
)

export default PageIntro
