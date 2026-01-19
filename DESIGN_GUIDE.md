# Idea Hamster - 90's Retro Design Guide 🎨

**Version:** 1.0  
**Last Updated:** January 2026  
**Purpose:** Visual design reference for consistent 90's aesthetic  

---

## 🎨 Color Palette

### Primary Colors
```css
/* Neon & Bright */
--neon-pink: #FF10F0;
--electric-blue: #00FFFF;
--lime-green: #39FF14;
--cyber-purple: #9D00FF;
--laser-yellow: #FFFF00;
--hot-magenta: #FF006E;

/* Backgrounds */
--deep-black: #0A0A0A;
--space-gray: #1A1A2E;
--midnight-blue: #16213E;

/* Accents */
--grid-blue: #0F3460;
--retro-orange: #FF6B35;
--arcade-red: #FF0055;
```

### Gradient Combinations
```css
/* Sunset Wave */
background: linear-gradient(135deg, #FF10F0 0%, #9D00FF 50%, #00FFFF 100%);

/* Cyber Grid */
background: linear-gradient(180deg, #0A0A0A 0%, #16213E 100%);

/* Neon Glow */
background: radial-gradient(circle, #FF10F0 0%, #0A0A0A 100%);

/* Retro Button */
background: linear-gradient(180deg, #9D00FF 0%, #6B00B3 100%);
```

---

## 🔤 Typography

### Font Stack

#### Headers & Titles
```css
/* Primary Display Font */
font-family: 'Press Start 2P', 'Courier New', monospace;
/* Use for: Logo, main headers, countdown timer */

/* Alternative Pixel Fonts */
font-family: 'VT323', monospace;  /* More readable for body text */
font-family: 'Orbitron', sans-serif;  /* Futuristic headers */
```

#### Body Text
```css
/* Monospace for Retro Feel */
font-family: 'Courier New', Courier, monospace;
font-family: 'Monaco', 'Lucida Console', monospace;
```

#### Special Elements
```css
/* Marquee Text */
font-family: 'Impact', 'Arial Black', sans-serif;
font-weight: bold;
text-transform: uppercase;

/* Terminal/Code */
font-family: 'IBM Plex Mono', 'Consolas', monospace;
```

### Type Scale
```css
/* Press Start 2P (smaller sizes for readability) */
--text-xs: 8px;   /* Tiny labels */
--text-sm: 10px;  /* Small UI text */
--text-base: 12px; /* Body text */
--text-lg: 16px;  /* Subheadings */
--text-xl: 20px;  /* Section headers */
--text-2xl: 28px; /* Page titles */
--text-3xl: 36px; /* Hero text */

/* VT323 (can go larger) */
--vt-base: 20px;
--vt-lg: 28px;
--vt-xl: 36px;
```

### Text Effects
```css
/* Neon Glow */
.neon-text {
  color: #00FFFF;
  text-shadow: 
    0 0 5px #00FFFF,
    0 0 10px #00FFFF,
    0 0 20px #00FFFF,
    0 0 40px #00FFFF;
}

/* 3D Retro */
.retro-3d {
  color: #FF10F0;
  text-shadow:
    2px 2px 0px #9D00FF,
    4px 4px 0px #6B00B3,
    6px 6px 0px #4A007A;
}

/* Pixel Shadow */
.pixel-shadow {
  color: #FFFF00;
  text-shadow:
    2px 2px 0px #000,
    2px -2px 0px #000,
    -2px 2px 0px #000,
    -2px -2px 0px #000;
}

/* Scanline Text */
.scanline-text {
  background: repeating-linear-gradient(
    0deg,
    rgba(0, 0, 0, 0.15),
    rgba(0, 0, 0, 0.15) 1px,
    transparent 1px,
    transparent 2px
  );
}
```

---

## 🎯 UI Components

### Buttons

#### Primary Button (Neon Style)
```css
.btn-primary {
  background: linear-gradient(180deg, #9D00FF 0%, #6B00B3 100%);
  border: 2px solid #FF10F0;
  color: #FFFFFF;
  font-family: 'Press Start 2P', monospace;
  font-size: 12px;
  padding: 12px 24px;
  text-transform: uppercase;
  cursor: pointer;
  box-shadow: 
    0 0 10px #FF10F0,
    inset 0 0 10px rgba(255, 16, 240, 0.3);
  transition: all 0.2s;
}

.btn-primary:hover {
  background: linear-gradient(180deg, #B300FF 0%, #8000D7 100%);
  box-shadow: 
    0 0 20px #FF10F0,
    0 0 40px #FF10F0,
    inset 0 0 10px rgba(255, 16, 240, 0.5);
  transform: translateY(-2px);
}

.btn-primary:active {
  transform: translateY(0);
  box-shadow: 
    0 0 5px #FF10F0,
    inset 0 0 10px rgba(255, 16, 240, 0.7);
}
```

#### Secondary Button (Retro 3D)
```css
.btn-secondary {
  background: #00FFFF;
  color: #0A0A0A;
  border: 3px solid #008B8B;
  font-family: 'Press Start 2P', monospace;
  font-size: 10px;
  padding: 10px 20px;
  position: relative;
  box-shadow: 
    3px 3px 0px #008B8B,
    6px 6px 0px #005555;
}

.btn-secondary:hover {
  top: 2px;
  left: 2px;
  box-shadow: 
    1px 1px 0px #008B8B,
    2px 2px 0px #005555;
}
```

#### Vote Button
```css
.btn-vote {
  background: linear-gradient(135deg, #39FF14 0%, #00CC00 100%);
  border: 2px solid #39FF14;
  color: #0A0A0A;
  font-family: 'VT323', monospace;
  font-size: 20px;
  padding: 8px 16px;
  position: relative;
  overflow: hidden;
}

.btn-vote::before {
  content: '▲';
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  font-size: 24px;
  opacity: 0;
  transition: opacity 0.3s;
}

.btn-vote:hover::before {
  opacity: 0.3;
}

.btn-vote.voted {
  background: linear-gradient(135deg, #666 0%, #444 100%);
  border-color: #666;
  color: #AAA;
  cursor: not-allowed;
}
```

### Cards

#### Idea Card
```css
.idea-card {
  background: #16213E;
  border: 2px solid #00FFFF;
  border-radius: 0; /* No rounded corners - pure 90's */
  padding: 16px;
  position: relative;
  box-shadow: 
    4px 4px 0px #0F3460,
    0 0 20px rgba(0, 255, 255, 0.3);
}

.idea-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: repeating-linear-gradient(
    0deg,
    transparent,
    transparent 2px,
    rgba(0, 255, 255, 0.03) 2px,
    rgba(0, 255, 255, 0.03) 4px
  );
  pointer-events: none;
}

.idea-card:hover {
  border-color: #FF10F0;
  box-shadow: 
    4px 4px 0px #0F3460,
    0 0 30px rgba(255, 16, 240, 0.5);
  transform: translateY(-2px);
}
```

#### Winner Card
```css
.winner-card {
  background: linear-gradient(135deg, #FFD700 0%, #FFA500 100%);
  border: 4px solid #FFFF00;
  padding: 24px;
  position: relative;
  animation: pulse-glow 2s infinite;
}

@keyframes pulse-glow {
  0%, 100% {
    box-shadow: 0 0 20px #FFFF00;
  }
  50% {
    box-shadow: 0 0 40px #FFFF00, 0 0 60px #FFA500;
  }
}

.winner-card::before {
  content: '★';
  position: absolute;
  top: -15px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 30px;
  color: #FFFF00;
  text-shadow: 0 0 10px #FFFF00;
}
```

### Leaderboard

```css
.leaderboard {
  background: #0A0A0A;
  border: 4px solid #00FFFF;
  padding: 20px;
  position: relative;
}

.leaderboard::before {
  content: '▓▓▓ LEADERBOARD ▓▓▓';
  position: absolute;
  top: -20px;
  left: 50%;
  transform: translateX(-50%);
  background: #0A0A0A;
  padding: 0 20px;
  font-family: 'Press Start 2P', monospace;
  font-size: 12px;
  color: #00FFFF;
  text-shadow: 0 0 10px #00FFFF;
}

.leaderboard-row {
  display: flex;
  align-items: center;
  padding: 12px;
  border-bottom: 1px solid #0F3460;
  transition: background 0.2s;
}

.leaderboard-row:hover {
  background: rgba(0, 255, 255, 0.1);
}

.leaderboard-rank {
  font-family: 'Press Start 2P', monospace;
  font-size: 16px;
  width: 40px;
  text-align: center;
}

.leaderboard-rank.first {
  color: #FFD700;
  text-shadow: 0 0 10px #FFD700;
}

.leaderboard-rank.second {
  color: #C0C0C0;
  text-shadow: 0 0 10px #C0C0C0;
}

.leaderboard-rank.third {
  color: #CD7F32;
  text-shadow: 0 0 10px #CD7F32;
}
```

### Progress Bars

#### Retro Loading Bar
```css
.progress-bar {
  width: 100%;
  height: 30px;
  background: #0A0A0A;
  border: 2px solid #00FFFF;
  position: relative;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: repeating-linear-gradient(
    90deg,
    #00FFFF 0px,
    #00FFFF 10px,
    #00CCCC 10px,
    #00CCCC 20px
  );
  transition: width 0.5s ease;
  position: relative;
  animation: shimmer 2s infinite;
}

@keyframes shimmer {
  0% {
    background-position: 0px 0px;
  }
  100% {
    background-position: 40px 0px;
  }
}

.progress-text {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-family: 'Press Start 2P', monospace;
  font-size: 10px;
  color: #FFFFFF;
  text-shadow: 1px 1px 0px #000;
  z-index: 10;
}
```

### Countdown Timer

```css
.countdown-timer {
  display: flex;
  gap: 20px;
  justify-content: center;
  padding: 20px;
  background: #0A0A0A;
  border: 3px solid #FF10F0;
  position: relative;
}

.countdown-segment {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 80px;
}

.countdown-number {
  font-family: 'Press Start 2P', monospace;
  font-size: 36px;
  color: #FF10F0;
  text-shadow: 
    0 0 10px #FF10F0,
    0 0 20px #FF10F0;
  background: #16213E;
  padding: 10px 20px;
  border: 2px solid #FF10F0;
  min-width: 100px;
  text-align: center;
}

.countdown-label {
  font-family: 'VT323', monospace;
  font-size: 16px;
  color: #00FFFF;
  margin-top: 8px;
  text-transform: uppercase;
}

.countdown-timer.urgent .countdown-number {
  animation: blink 1s infinite;
  color: #FF0055;
  border-color: #FF0055;
  text-shadow: 
    0 0 10px #FF0055,
    0 0 20px #FF0055;
}

@keyframes blink {
  0%, 49% { opacity: 1; }
  50%, 100% { opacity: 0.5; }
}
```

---

## 🎭 Special Effects

### CRT Screen Effect
```css
.crt-effect {
  position: relative;
  overflow: hidden;
}

.crt-effect::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: repeating-linear-gradient(
    0deg,
    rgba(0, 0, 0, 0.15),
    rgba(0, 0, 0, 0.15) 1px,
    transparent 1px,
    transparent 3px
  );
  pointer-events: none;
  z-index: 1000;
}

.crt-effect::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(18, 16, 16, 0.1);
  pointer-events: none;
  animation: flicker 0.15s infinite;
}

@keyframes flicker {
  0% { opacity: 0.27861; }
  5% { opacity: 0.34769; }
  10% { opacity: 0.23604; }
  15% { opacity: 0.90626; }
  20% { opacity: 0.18128; }
  25% { opacity: 0.83891; }
  30% { opacity: 0.65583; }
  35% { opacity: 0.67807; }
  40% { opacity: 0.26559; }
  45% { opacity: 0.84693; }
  50% { opacity: 0.96019; }
  55% { opacity: 0.08594; }
  60% { opacity: 0.20313; }
  65% { opacity: 0.71988; }
  70% { opacity: 0.53455; }
  75% { opacity: 0.37288; }
  80% { opacity: 0.71428; }
  85% { opacity: 0.70419; }
  90% { opacity: 0.7003; }
  95% { opacity: 0.36108; }
  100% { opacity: 0.24387; }
}
```

### Glitch Effect (on Click)
```css
.glitch {
  position: relative;
}

.glitch.active {
  animation: glitch-skew 0.3s;
}

@keyframes glitch-skew {
  0% { transform: skew(0deg); }
  10% { transform: skew(-5deg); }
  20% { transform: skew(5deg); }
  30% { transform: skew(-5deg); }
  40% { transform: skew(5deg); }
  50% { transform: skew(0deg); }
  100% { transform: skew(0deg); }
}

.glitch.active::before {
  content: attr(data-text);
  position: absolute;
  left: 2px;
  text-shadow: -2px 0 #FF10F0;
  animation: glitch-anim-1 0.3s;
}

.glitch.active::after {
  content: attr(data-text);
  position: absolute;
  left: -2px;
  text-shadow: 2px 0 #00FFFF;
  animation: glitch-anim-2 0.3s;
}

@keyframes glitch-anim-1 {
  0% { clip-path: inset(40% 0 61% 0); }
  20% { clip-path: inset(92% 0 1% 0); }
  40% { clip-path: inset(43% 0 1% 0); }
  60% { clip-path: inset(25% 0 58% 0); }
  80% { clip-path: inset(54% 0 7% 0); }
  100% { clip-path: inset(58% 0 43% 0); }
}
```

### Pixel Transition
```css
@keyframes pixel-wipe {
  0% {
    clip-path: polygon(0 0, 0 0, 0 100%, 0 100%);
  }
  100% {
    clip-path: polygon(0 0, 100% 0, 100% 100%, 0 100%);
  }
}

.pixel-transition {
  animation: pixel-wipe 0.5s steps(8) forwards;
}
```

### Starfield Background
```css
.starfield {
  background: #0A0A0A;
  position: relative;
  overflow: hidden;
}

.star {
  position: absolute;
  width: 2px;
  height: 2px;
  background: white;
  border-radius: 50%;
  animation: twinkle 3s infinite;
}

@keyframes twinkle {
  0%, 100% { opacity: 0.3; }
  50% { opacity: 1; }
}

/* Generate stars with JavaScript */
```

---

## 🎮 Animated Elements

### Hamster Mascot States

#### Idle Animation
```css
.hamster-idle {
  width: 64px;
  height: 64px;
  background-image: url('/hamster-idle-sprite.png');
  animation: idle-breath 2s infinite;
}

@keyframes idle-breath {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.05); }
}
```

#### Running Animation
```css
.hamster-running {
  width: 64px;
  height: 64px;
  background-image: url('/hamster-run-sprite.png');
  animation: run-cycle 0.6s steps(8) infinite;
}

@keyframes run-cycle {
  0% { background-position: 0px 0px; }
  100% { background-position: -512px 0px; }
}
```

#### Celebrating Animation
```css
.hamster-celebrate {
  width: 64px;
  height: 64px;
  background-image: url('/hamster-celebrate-sprite.png');
  animation: celebrate 0.8s steps(10) infinite;
}

@keyframes celebrate {
  0% { background-position: 0px 0px; }
  100% { background-position: -640px 0px; }
}
```

### Loading Spinner
```css
.retro-spinner {
  width: 50px;
  height: 50px;
  border: 4px solid #0F3460;
  border-top-color: #00FFFF;
  border-radius: 0; /* Keep it square for retro feel */
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
```

### Marquee Text
```css
.marquee {
  width: 100%;
  overflow: hidden;
  background: #0A0A0A;
  border-top: 2px solid #FF10F0;
  border-bottom: 2px solid #FF10F0;
  padding: 8px 0;
}

.marquee-content {
  display: inline-block;
  white-space: nowrap;
  animation: marquee-scroll 20s linear infinite;
  font-family: 'Press Start 2P', monospace;
  font-size: 12px;
  color: #00FFFF;
}

@keyframes marquee-scroll {
  0% { transform: translateX(100%); }
  100% { transform: translateX(-100%); }
}
```

---

## 🖼️ Image Styles

### Pixel Art Images
```css
.pixel-image {
  image-rendering: pixelated;
  image-rendering: -moz-crisp-edges;
  image-rendering: crisp-edges;
}
```

### Polaroid Photo Style
```css
.polaroid {
  background: white;
  padding: 10px;
  padding-bottom: 40px;
  box-shadow: 
    0 4px 6px rgba(0, 0, 0, 0.3),
    0 0 20px rgba(0, 255, 255, 0.2);
  transform: rotate(-2deg);
  transition: transform 0.3s;
}

.polaroid:hover {
  transform: rotate(0deg) scale(1.05);
}

.polaroid img {
  width: 100%;
  display: block;
  border: 1px solid #eee;
}
```

---

## 🎵 Sound Design

### Sound Effects Needed
- **Vote Cast:** `beep.mp3` (short, uplifting)
- **Cycle End:** `victory.mp3` (8-bit fanfare)
- **Build Start:** `startup.mp3` (computer boot sound)
- **Build Complete:** `level-up.mp3` (achievement sound)
- **Error:** `error.mp3` (bonk/buzz)
- **Navigation:** `click.mp3` (subtle button press)

### Volume Guidelines
```javascript
const SOUND_VOLUMES = {
  votecast: 0.3,
  cycleEnd: 0.5,
  buildStart: 0.4,
  buildComplete: 0.6,
  error: 0.4,
  navigation: 0.2
};
```

---

## 📱 Responsive Design

### Breakpoints
```css
/* Mobile First Approach */
--mobile: 320px;
--tablet: 768px;
--desktop: 1024px;
--wide: 1440px;

/* Adjust pixel font sizes for mobile */
@media (max-width: 768px) {
  .countdown-number {
    font-size: 24px;
  }
  
  .leaderboard-rank {
    font-size: 12px;
  }
  
  .btn-primary {
    font-size: 10px;
    padding: 10px 16px;
  }
}
```

---

## ✨ Accessibility

### High Contrast Mode
```css
@media (prefers-contrast: high) {
  :root {
    --neon-pink: #FF00FF;
    --electric-blue: #00FFFF;
    --lime-green: #00FF00;
  }
  
  .neon-text {
    text-shadow: none;
    outline: 2px solid currentColor;
  }
}
```

### Reduced Motion
```css
@media (prefers-reduced-motion: reduce) {
  * {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
  
  .crt-effect::after {
    animation: none;
  }
}
```

---

## 🎯 Component Library

### Recommended Structure
```
components/
├── retro/
│   ├── Button.tsx
│   ├── Card.tsx
│   ├── ProgressBar.tsx
│   ├── CountdownTimer.tsx
│   ├── Leaderboard.tsx
│   ├── Marquee.tsx
│   ├── HamsterSprite.tsx
│   └── GlitchText.tsx
├── layout/
│   ├── Header.tsx
│   ├── Footer.tsx
│   └── CRTWrapper.tsx
└── effects/
    ├── Scanlines.tsx
    ├── Starfield.tsx
    └── PixelTransition.tsx
```

---

## 🎨 Design Resources

### Free Pixel Art Tools
- **Piskel** - Browser-based sprite editor
- **Aseprite** - Professional pixel art tool
- **GIMP** - With pixel art plugins

### Font Sources
- **Google Fonts:** Press Start 2P, VT323, Orbitron
- **DaFont:** Retro gaming fonts section
- **Font Squirrel:** Free pixel fonts

### Inspiration Sites
- **Poolsuite.net** - Excellent retro web design
- **CodePen** - Search "retro" or "vaporwave"
- **Dribbble** - 90's UI design tags

### Color Tools
- **Coolors.co** - Neon palette generator
- **Adobe Color** - Create retro schemes
- **Vaporwave Color Palette Generator**

---

## 🚀 Implementation Tips

1. **Start with the grid:** Establish your layout using CSS Grid
2. **Layer effects:** Scanlines and CRT effects should be subtle
3. **Test on real devices:** Pixel fonts can be hard to read on small screens
4. **Optimize performance:** Limit animations, use CSS transforms
5. **Provide off-switches:** Let users disable effects if needed
6. **Balance nostalgia and usability:** Don't sacrifice UX for aesthetics

---

**Remember:** The goal is "Fun 90's" not "Unusable 90's"! 🎮✨