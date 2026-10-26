const PageIntro = ({ eyebrow, title, children, meta, constellation = false }) => (
  <header className={`page-intro${constellation ? ' page-intro--constellation' : ''}`}>
    {constellation && <div className="constellation" aria-hidden="true" />}
    <div className="page-intro__content">
      <p className="eyebrow">{eyebrow}</p>
      <h1>{title}</h1>
      {children && <p className="page-intro__copy">{children}</p>}
      {meta && <p className="page-intro__meta">{meta}</p>}
    </div>
  </header>
)

export default PageIntro
