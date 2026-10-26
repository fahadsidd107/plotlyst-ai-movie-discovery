import icon from '../../assets/PlotlystIcon.png';

const LoadingScreen = ({ compact = false }) => (
  <div className={`plotlyst-loader${compact ? ' plotlyst-loader--compact' : ''}`} role="status" aria-live="polite">
    <div className="plotlyst-loader__orbit" aria-hidden="true">
      <span /><span /><span />
      <img src={icon} alt="" />
    </div>
    <p className="eyebrow">Plotlyst is cueing your next obsession</p>
    <strong>Stories are aligning.</strong>
    <span className="visually-hidden">Loading Plotlyst</span>
  </div>
);

export default LoadingScreen;
