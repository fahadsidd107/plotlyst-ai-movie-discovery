import icon from '../../assets/PlotlystIcon.png';

const LoadingScreen = ({ compact = false }) => (
  <div className={`plotlyst-loader${compact ? ' plotlyst-loader--compact' : ''}`} role="status" aria-live="polite">
    <div className="plotlyst-loader__gate" aria-hidden="true">
      <span className="plotlyst-loader__count">3</span>
      <span className="plotlyst-loader__sprockets" />
      <img src={icon} alt="" />
    </div>
    <p className="eyebrow">Projector warming · Reel one</p>
    <strong>The programme starts shortly.</strong>
    <span className="visually-hidden">Loading Plotlyst</span>
  </div>
);

export default LoadingScreen;
