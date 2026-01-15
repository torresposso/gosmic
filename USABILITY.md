# 🛸 USABILITY.md — Gosmic UX & Accessibility Audit

> **Fecha:** 2026-01-14  
> **Versión:** v1.0.0  
> **Auditor:** Frontend Engineer (Space Enthusiast 🚀)  
> **Stack Analizado:** DaisyUI v5 + Tailwind CSS v4 + Alpine.js + Templ

---

## 📋 ÍNDICE

1. [Resumen Ejecutivo](#-resumen-ejecutivo)
2. [Usabilidad (UX)](#-usabilidad-ux)
3. [Accesibilidad (A11Y)](#-accesibilidad-a11y)
4. [Diseño Visual](#-diseño-visual)
5. [Rendimiento Percibido](#-rendimiento-percibido)
6. [Recomendaciones Priorizadas](#-recomendaciones-priorizadas)
7. [Checklist WCAG 2.1](#-checklist-wcag-21)

---

## 📊 RESUMEN EJECUTIVO

| Área | Estado | Puntuación |
|------|--------|------------|
| **Usabilidad General** | ⚠️ Buena | 7.5/10 |
| **Accesibilidad (WCAG 2.1)** | 🔴 Necesita Mejoras | 5.5/10 |
| **Diseño Visual** | ✅ Excelente | 9/10 |
| **Rendimiento Percibido** | ✅ Muy Bueno | 8.5/10 |
| **Navegación** | ⚠️ Buena | 7/10 |

**Puntuación General UX/A11Y:** 7.3/10 ⭐

### 🎯 Hallazgos Críticos

1. **🔴 CRÍTICO:** Los SVGs carecen de atributos `aria-label` o texto alternativo
2. **🔴 CRÍTICO:** Falta atributo `role="main"` en el contenedor principal
3. **🟡 ALTO:** Toggle switches sin labels accesibles asociados
4. **🟡 ALTO:** Los flash messages desaparecen automáticamente sin control del usuario
5. **🟢 BUENO:** Excelente uso de `lang="en"` y viewport meta

---

## 🎨 USABILIDAD (UX)

### ✅ Puntos Fuertes — "Houston, tenemos éxito"

| Aspecto | Descripción |
|---------|-------------|
| **Consistencia Visual** | Tema espacial cohesivo en toda la app. Cada elemento "encaja" como módulos de una estación espacial. |
| **Navegación Clara** | Menú desktop e móvil bien diferenciados. Los enlaces son intuitivos. |
| **Feedback Visual** | Excelentes transiciones hover y estados focus con efectos de "hyperdrive". |
| **Responsividad** | Diseño mobile-first con breakpoints apropiados (`lg:hidden`, `lg:flex`). |
| **Temática Inmersiva** | Terminología espacial (Commander, Mission Logs, Abort Session) crea una experiencia memorable. |

### ⚠️ Áreas de Mejora — "Control de misión reporta anomalías"

#### 1. Formularios — Friction de Entrada

```
📍 Ubicación: home.templ (líneas 185-222), posts.templ (líneas 41-84), auth.templ
```

| Problema | Impacto | Severidad |
|----------|---------|-----------|
| Placeholders como únicos indicadores | Los placeholders desaparecen al escribir, el usuario pierde contexto | 🟡 Medio |
| Sin validación en tiempo real | El usuario descubre errores solo al enviar | 🟡 Medio |
| Toggle "Public" sin label visible | Solo hay texto técnico "Deep_Space_Broadcast" | 🟡 Medio |

**Recomendación Cósmica:** Implementar labels persistentes flotantes ("floating labels") y mensajes de validación inline.

#### 2. Flash Messages — Órbita Inestable

```
📍 Ubicación: layout.templ (líneas 87-113)
```

| Problema | Impacto | Severidad |
|----------|---------|-----------|
| Auto-dismissal a 5s sin pausa | Usuarios con discapacidades cognitivas pueden perder el mensaje | 🟡 Medio |
| No hay botón de cierre | El usuario no tiene control sobre notificaciones | 🟢 Bajo |
| Posición fixed en móvil | Puede superponerse a contenido importante | 🟢 Bajo |

#### 3. Delete Confirmation — Peligro sin Escudo

```
📍 Ubicación: posts.templ (líneas 138-147)
```

| Problema | Impacto | Severidad |
|----------|---------|-----------|
| Acción destructiva sin confirmación | Un click accidental elimina datos permanentemente | 🔴 Alto |
| Botón DELETE muy cerca de EDIT | Facilita errores de "fat finger" en móvil | 🟡 Medio |

**Recomendación Cósmica:** Implementar modal de confirmación con Alpine.js:
```html
<button @click="confirmDelete = true">Purge</button>
<dialog x-show="confirmDelete">¿Confirmar expulsión al vacío?</dialog>
```

#### 4. Navegación — Waypoints Faltantes

| Problema | Impacto | Severidad |
|----------|---------|-----------|
| Sin breadcrumbs en páginas internas | El usuario pierde noción de ubicación | 🟢 Bajo |
| Sin indicador de página activa en navbar | No hay feedback visual del estado actual | 🟡 Medio |
| Logo no tiene texto "Home" | Usuarios screen reader no saben que lleva al inicio | 🟡 Medio |

---

## ♿ ACCESIBILIDAD (A11Y)

### 🔴 Problemas Críticos — "Alerta Roja en el Puente"

#### 1. SVGs Sin Texto Alternativo

```
📍 Ubicación: TODOS los templates
```

**Problema:** Cada ícono SVG carece de `aria-label` o `<title>`:

```html
<!-- ❌ ACTUAL -->
<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
  <path .../>
</svg>

<!-- ✅ RECOMENDADO -->
<svg aria-label="Icono de búsqueda" role="img" ...>
  <title>Buscar</title>
  <path .../>
</svg>
```

**Conteo de SVGs sin a11y:** 25+ instancias 🚨

#### 2. Landmark Roles Faltantes

```
📍 Ubicación: layout.templ
```

| Elemento | Actual | Requerido |
|----------|--------|-----------|
| `<nav>` | `<div class="navbar">` | `<nav aria-label="Main navigation">` |
| `<main>` | `<main class="...">` | ✅ Correcto |
| `<footer>` | `<footer>` | ✅ Correcto |
| Skip Link | ❌ No existe | `<a href="#main" class="skip-link">Skip to content</a>` |

#### 3. Labels de Formulario Incompletos

```
📍 Ubicación: home.templ, posts.templ
```

**Toggle Switch sin asociación:**
```html
<!-- ❌ ACTUAL: El label y el input no están asociados -->
<span class="label-text">Deep_Space_Broadcast</span>
<input type="checkbox" id="public" name="public" class="toggle toggle-primary"/>

<!-- ✅ RECOMENDADO -->
<label for="public" class="flex items-center gap-2">
  <span>Deep Space Broadcast (Público)</span>
  <input type="checkbox" id="public" name="public" class="toggle toggle-primary"/>
</label>
```

#### 4. Focus Management

| Problema | Ubicación | Severidad |
|----------|-----------|-----------|
| Dropdown móvil sin focus trap | layout.templ:18-35 | 🟡 Medio |
| Flash toast no recibe focus | layout.templ:89 | 🟡 Medio |
| Modal edit sin focus inicial | posts.templ:153-197 | 🟡 Medio |

### 🟡 Problemas Moderados

#### 5. Contraste de Color

| Elemento | Ratio Actual | WCAG AA | Veredicto |
|----------|--------------|---------|-----------|
| `text-primary/60` labels | ~4.2:1 | 4.5:1 | ⚠️ Borderline |
| `text-primary/40` hints | ~2.8:1 | 4.5:1 | 🔴 Falla |
| `opacity-50` footer | ~3.1:1 | 4.5:1 | 🔴 Falla |

**Afectados:**
- Labels de formulario con `text-primary/60`
- Texto "Terminal_ID", "Status:" en dashboard
- Footer "Powered by..."

#### 6. Emojis Sin Descripción

Los emojis usados como íconos deben envolverse:
```html
<!-- ❌ ACTUAL -->
<span class="text-primary">🚀</span>

<!-- ✅ RECOMENDADO -->
<span role="img" aria-label="Cohete">🚀</span>
```

---

## 🎨 DISEÑO VISUAL

### ✅ Excelencia Galáctica

| Aspecto | Evaluación |
|---------|------------|
| **Paleta de Colores** | Tema "night" con acentos primary perfectamente calibrados. Recuerda a los paneles de control de la ISS. |
| **Tipografía** | Uso de `font-mono` para elementos técnicos crea estética cyberpunk/espacial. |
| **Espaciado** | Grid system con gaps consistentes (`gap-6`, `gap-4`). |
| **Iconografía** | SVGs inline con stroke-currentColor permiten theming coherente. |
| **Micro-interacciones** | Efectos hover con `translate-y-1` y sombras dinámicas (`hover:shadow-primary/20`). |
| **HUD Corners** | Los bordes decorativos simulan interfaz de nave espacial. Muy inmersivo. |

### 🌟 Elementos Destacados

1. **Backdrop Blur Navbar:** El efecto glassmorphism (`backdrop-blur-md`) es moderno y funcional.
2. **Glow Effects:** Los botones con `shadow-[0_0_20px...]` crean efecto neón espacial.
3. **Animate Pulse:** El ícono de transmisión pulsante refuerza la temática de broadcast.

---

## ⚡ RENDIMIENTO PERCIBIDO

### ✅ Optimizaciones Detectadas

| Técnica | Implementación |
|---------|----------------|
| **Alpine.js defer** | `<script src="..." defer>` — Carga no bloqueante |
| **CSS Minimal** | DaisyUI tree-shaking via `@source` directive |
| **Transiciones GPU** | Uso de `transform` y `opacity` |
| **Iconos Inline** | SVGs embebidos evitan peticiones adicionales |

### ⚠️ Oportunidades

| Aspecto | Recomendación |
|---------|---------------|
| **Preload Fonts** | Añadir `<link rel="preload">` para fuentes mono |
| **Loading States** | Añadir skeleton placeholders mientras cargan posts |
| **Image Optimization** | No hay imágenes pesadas, pero considerar lazy loading futuro |

---

## 🎯 RECOMENDACIONES PRIORIZADAS

### 🔴 Prioridad Alta — "Ejecutar Inmediatamente"

| # | Tarea | Archivo | Esfuerzo |
|---|-------|---------|----------|
| 1 | Añadir `aria-label` a todos los SVGs | Todos | 1h |
| 2 | Convertir `.navbar` a `<nav>` semántico | layout.templ | 10min |
| 3 | Añadir skip link "Skip to content" | layout.templ | 15min |
| 4 | Implementar confirmación para DELETE | posts.templ | 30min |
| 5 | Asociar labels con toggles (for/id) | home.templ, posts.templ | 20min |

### 🟡 Prioridad Media — "Siguiente Sprint"

| # | Tarea | Archivo | Esfuerzo |
|---|-------|---------|----------|
| 6 | Aumentar contraste de textos `/60` y `/40` | Todos | 30min |
| 7 | Añadir botón dismiss a flash messages | layout.templ | 20min |
| 8 | Implementar focus trap en dropdown móvil | layout.templ | 45min |
| 9 | Añadir indicador de página activa en navbar | layout.templ | 30min |
| 10 | Envolver emojis en `role="img"` | Todos | 30min |

### 🟢 Prioridad Baja — "Nice to Have"

| # | Tarea | Descripción |
|---|-------|-------------|
| 11 | Breadcrumbs | Añadir migas de pan en dashboard/posts |
| 12 | Loading skeletons | Placeholders mientras cargan posts |
| 13 | Dark/Light toggle | Permitir cambio de tema (ya soportado por DaisyUI) |
| 14 | Animaciones reducidas | Respetar `prefers-reduced-motion` |

---

## ✅ CHECKLIST WCAG 2.1 (Nivel AA)

### Perceptible

| Criterio | Estado | Notas |
|----------|--------|-------|
| 1.1.1 Non-text Content | 🔴 Falla | SVGs sin alt |
| 1.3.1 Info and Relationships | 🟡 Parcial | Faltan landmarks semánticos |
| 1.3.2 Meaningful Sequence | ✅ Pasa | DOM order correcto |
| 1.4.3 Contrast (Minimum) | 🔴 Falla | Textos con opacidad baja |
| 1.4.4 Resize Text | ✅ Pasa | Responsive scaling |

### Operable

| Criterio | Estado | Notas |
|----------|--------|-------|
| 2.1.1 Keyboard | 🟡 Parcial | Focus visible pero no hay skip link |
| 2.4.1 Bypass Blocks | 🔴 Falla | Sin skip link |
| 2.4.2 Page Titled | ✅ Pasa | Títulos descriptivos |
| 2.4.6 Headings and Labels | 🟡 Parcial | H1-H2 correctos, labels incompletos |

### Comprensible

| Criterio | Estado | Notas |
|----------|--------|-------|
| 3.1.1 Language of Page | ✅ Pasa | `lang="en"` presente |
| 3.2.2 On Input | ✅ Pasa | Sin cambios inesperados |
| 3.3.1 Error Identification | ✅ Pasa | Alertas de error visibles |
| 3.3.2 Labels or Instructions | 🟡 Parcial | Algunos campos solo con placeholder |

### Robusto

| Criterio | Estado | Notas |
|----------|--------|-------|
| 4.1.1 Parsing | ✅ Pasa | HTML válido |
| 4.1.2 Name, Role, Value | 🔴 Falla | Controles sin nombres accesibles |

---

## 🚀 CONCLUSIÓN

Gosmic tiene una **identidad visual excepcional** que captura perfectamente la esencia de una aplicación espacial. La temática es consistente y las micro-interacciones elevaan la experiencia.

Sin embargo, la accesibilidad necesita atención urgente. Con ~3 horas de trabajo enfocado en los items de prioridad alta, podríamos elevar la puntuación A11Y de 5.5/10 a 8/10, abriendo la app a un universo más amplio de usuarios.

**¡Ad Astra Per Aspera, Commander!** ✨

---

*Documento generado automáticamente. Última actualización: 2026-01-14 21:38 EST*
