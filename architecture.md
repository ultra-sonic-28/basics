# Architecture "Basics"

## Overview

The project is a multi-BASIC, multi-architecture interpreter designed for retro computers (Apple II, C64, etc.). It uses a layered architecture to separate the BASIC language logic from the hardware-specific emulation and the graphical rendering.

## Component Diagram

```mermaid
flowchart TB

    %% =====================
    %% BASIC CORE
    %% =====================
    subgraph BASIC_CORE["BASIC Core"]
        LEXER["lexer
        (source → tokens)"]
        PARSER["parser
        (tokens → AST)"]
        INTERPRETER["interpreter
        (AST → Runtime calls)
        Renderer agnostic"]
    end

    LEXER --> PARSER
    PARSER --> INTERPRETER

    %% =====================
    %% RUNTIME & ENVIRONMENT
    %% =====================
    subgraph RUNTIME_LAYER["Runtime Layer"]
        RT["Runtime
        (BASIC API: PRINT, HOME, PLOT, PR#...)"]
        ENV["Environment
        (Variables, Arrays, Stacks)"]
    end

    BASIC_CORE --> RT
    RT --> ENV

    %% =====================
    %% HARDWARE ABSTRACTION
    %% =====================
    subgraph HARDWARE_LAYER["Hardware Abstraction (video)"]
        DEVICE_IF["video.Device (Interface)
        Logic seen by Runtime"]
        
        subgraph APPLE2_DM["Apple II DisplayManager"]
            DM["DisplayManager
            (Manages multiple modes)"]
            MODE_IF["video.Mode (Interface)"]
            APPLE_TEXT["AppleText
            (Unified 40/80 cols)"]
        end
        
        TTY["TTY Device
        (Terminal output)"]
    end

    RT --> DEVICE_IF
    DEVICE_IF --> DM
    DEVICE_IF --> TTY
    DM --> MODE_IF
    MODE_IF --> APPLE_TEXT

    %% =====================
    %% RENDERING ENGINE
    %% =====================
    subgraph RENDERING_LAYER["Rendering Engine"]
        V_RENDERER["video.Renderer (Interface)
        DrawPixel, DrawGlyph"]
        EBITEN_REND["Ebiten Renderer
        (Bitmap Framebuffer)"]
        FONT_MGR["Font Manager
        (Bitmap Fonts 7x8, 8x8)"]
    end

    APPLE_TEXT --> V_RENDERER
    V_RENDERER --> EBITEN_REND
    V_RENDERER --> FONT_MGR

    %% =====================
    %% EXECUTION HOST
    %% =====================
    EBITEN_APP["EbitenApp
    (Window Manager, Input, Resize)"]
    EBITEN_APP --> RT
    EBITEN_REND --> EBITEN_APP
```

## Key Architectural Concepts

### 1. DisplayManager & Multi-Mode Support
The `DisplayManager` (e.g., in `internal/machines/apple2`) acts as a proxy for the `video.Device` interface. It manages a collection of `video.Mode` instances (like `AppleText` or future `GR`/`HGR` modes).
- **Dynamic Switching**: The `SwitchMode(slot)` method (called via `PR#n`) allows the system to change the active mode at runtime.
- **Proxy Pattern**: All calls to the `video.Device` interface (Clear, Print, etc.) are delegated to the currently active `video.Mode`.

### 2. AppleText (Unified Text Mode)
Replacing the older `Text40` and `Text80` files, `AppleText` is a generic text mode implementation that:
- Uses `internal/video/text.TextMode` for buffer logic.
- Handles keyboard input (`INPUT`, `GET`) and blinking cursor logic.
- **Buffer Preservation**: Implements `CopyFrom(other)`, allowing content to be transferred when switching between 40 and 80 columns (or vice-versa), preserving the user's display.

### 3. Rendering Abstraction
The `video.Renderer` interface provides low-level primitives (`DrawPixel`, `DrawGlyph`). 
- **Ebiten Implementation**: The actual rendering is done via a bitmap framebuffer in `internal/video/ebiten`, which is then blitted to the screen.
- **Dynamic Resizing**: `EbitenApp` monitors the active mode's dimensions and automatically resizes the host window when a mode switch occurs (e.g., from 280px to 560px width).

### 4. Input Handling
Input is handled at the `EbitenApp` level and injected into the active `video.Device`. 
- **GET/INPUT Modes**: The device tracks whether it's waiting for a single key (`GET`) or a full line (`INPUT`).
- **Blink Cycle**: The `Update()` method in the device is called 60 times per second to handle internal state like cursor blinking.
