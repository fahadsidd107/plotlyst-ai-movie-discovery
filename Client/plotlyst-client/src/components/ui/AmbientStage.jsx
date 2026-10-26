import { useEffect, useRef } from 'react';

const AmbientStage = () => {
  const canvasRef = useRef(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    const context = canvas.getContext('2d');
    const reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    const pointer = { x: -1000, y: -1000 };
    const ripples = [];
    let frameId;
    let particles = [];

    const resize = () => {
      const ratio = Math.min(window.devicePixelRatio || 1, 2);
      canvas.width = window.innerWidth * ratio;
      canvas.height = window.innerHeight * ratio;
      canvas.style.width = `${window.innerWidth}px`;
      canvas.style.height = `${window.innerHeight}px`;
      context.setTransform(ratio, 0, 0, ratio, 0, 0);
      particles = Array.from({ length: Math.min(56, Math.floor(window.innerWidth / 24)) }, () => ({
        x: Math.random() * window.innerWidth,
        y: Math.random() * window.innerHeight,
        vx: (Math.random() - 0.5) * 0.18,
        vy: (Math.random() - 0.5) * 0.18,
        radius: Math.random() * 1.4 + 0.5,
        color: ['#a977ff', '#ff657a', '#6fe7da'][Math.floor(Math.random() * 3)],
      }));
    };

    const onPointerMove = (event) => {
      pointer.x = event.clientX;
      pointer.y = event.clientY;
    };
    const onPointerLeave = () => Object.assign(pointer, { x: -1000, y: -1000 });
    const onPointerDown = (event) => ripples.push({ x: event.clientX, y: event.clientY, radius: 0, alpha: 0.45 });

    const draw = () => {
      context.clearRect(0, 0, window.innerWidth, window.innerHeight);
      particles.forEach((particle, index) => {
        if (!reducedMotion) {
          particle.x += particle.vx;
          particle.y += particle.vy;
          if (particle.x < 0 || particle.x > window.innerWidth) particle.vx *= -1;
          if (particle.y < 0 || particle.y > window.innerHeight) particle.vy *= -1;
        }

        const distance = Math.hypot(pointer.x - particle.x, pointer.y - particle.y);
        if (!reducedMotion && distance < 130 && distance > 0) {
          particle.x -= ((pointer.x - particle.x) / distance) * 0.35;
          particle.y -= ((pointer.y - particle.y) / distance) * 0.35;
        }

        context.globalAlpha = distance < 150 ? 0.7 : 0.24;
        context.fillStyle = particle.color;
        context.beginPath();
        context.arc(particle.x, particle.y, particle.radius, 0, Math.PI * 2);
        context.fill();

        for (let otherIndex = index + 1; otherIndex < particles.length; otherIndex += 1) {
          const other = particles[otherIndex];
          const gap = Math.hypot(particle.x - other.x, particle.y - other.y);
          if (gap < 95) {
            context.globalAlpha = (1 - gap / 95) * 0.08;
            context.strokeStyle = particle.color;
            context.beginPath();
            context.moveTo(particle.x, particle.y);
            context.lineTo(other.x, other.y);
            context.stroke();
          }
        }
      });

      for (let index = ripples.length - 1; index >= 0; index -= 1) {
        const ripple = ripples[index];
        ripple.radius += reducedMotion ? 10 : 2.8;
        ripple.alpha -= reducedMotion ? 0.12 : 0.012;
        context.globalAlpha = Math.max(ripple.alpha, 0);
        context.strokeStyle = '#a977ff';
        context.lineWidth = 1.5;
        context.beginPath();
        context.arc(ripple.x, ripple.y, ripple.radius, 0, Math.PI * 2);
        context.stroke();
        if (ripple.alpha <= 0) ripples.splice(index, 1);
      }
      context.globalAlpha = 1;
      frameId = requestAnimationFrame(draw);
    };

    resize();
    draw();
    window.addEventListener('resize', resize);
    window.addEventListener('pointermove', onPointerMove);
    document.addEventListener('pointerleave', onPointerLeave);
    window.addEventListener('pointerdown', onPointerDown);
    return () => {
      cancelAnimationFrame(frameId);
      window.removeEventListener('resize', resize);
      window.removeEventListener('pointermove', onPointerMove);
      document.removeEventListener('pointerleave', onPointerLeave);
      window.removeEventListener('pointerdown', onPointerDown);
    };
  }, []);

  return <canvas ref={canvasRef} className="ambient-stage" aria-hidden="true" />;
};

export default AmbientStage;
