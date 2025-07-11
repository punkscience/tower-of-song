# **Tower of Song UI Specification**

This document outlines the user interface for the "Tower of Song" personal music streaming service, designed for a minimal, elegant, and thematically relevant aesthetic. It aims for responsiveness across web and mobile platforms.

## **I. Global Styles & Theme**

* **Color Palette:**  
  * **Background:** Deep charcoal or very dark blue (\#1A1A1A or similar dark tone).  
  * **Text (Primary):** Off-white or very light grey (\#F5F5F5).  
  * **Accent 1 (Interactive Elements, Highlights):** Muted gold or light bronze (\#D4AF37 or \#B8860B).  
  * **Accent 2 (Secondary Text, Subtle Borders):** Medium grey (\#888888).  
* **Typography:**  
  * **Primary Font (Body, Controls):** Clean sans-serif (e.g., Inter, Roboto, or Lato \- for legibility).  
  * **Secondary Font (Title, Headings):** A slightly more elegant sans-serif or a subtle serif if desired for Tower of Song title (e.g., Lora, Merriweather, or slightly bolder weight of primary font).  
* **Spacing:** Generous use of negative space (padding, margins) to ensure a clean, uncluttered look.  
* **Borders/Separators:** Minimal, subtle thin lines if needed, using Accent 2\.

## **II. Layout Structure**

The layout will adapt using **flexbox and grid** for responsiveness.

* **Mobile (Portrait):** Stacked vertical sections.  
* **Web/Tablet (Landscape):** Utilizes horizontal space, potentially side-by-side arrangement for controls and lists.

### **A. Top Section: Header & Now Playing**

* **Container:** div with flex properties, potentially justify-between on web, flex-col or items-center on mobile.  
* **Elements:**  
  1. **"Tower of Song" Title:**  
     * **Tag:** \<h1\> or \<div\> for branding.  
     * **Text:** "Tower of Song".  
     * **Styling:** text-3xl to text-4xl on web, text-2xl to text-3xl on mobile. font-bold or specific font-serif if chosen. text-accent1 or text-primary.  
     * **Placement:** top-left on web, center-top on mobile.  
  2. **Currently Playing Information:**  
     * **Container:** div to group artist, album, track, and favorite toggle.  
     * **Artist:** text-lg to text-xl, text-primary.  
     * **Album:** text-md to text-lg, text-accent2.  
     * **Track Title:** text-xl to text-2xl (font-bold), text-primary.  
     * **"Favourite" Toggle:**  
       * **Tag:** \<button\> or \<div\> containing an icon.  
       * **Icon:** ♡ (unfavorited) / ❤️ (favorited). Can be SVG or icon font.  
       * **Styling:** text-primary when unfavorited, text-red-500 or text-accent1 when favorited. cursor-pointer.  
       * **Placement:** Directly to the right of the track title.  
       * **Interaction:** Toggles favorite status on click.

### **B. Middle Section: Controls & Search**

* **Container:** div to group player controls and search functionality.  
1. **Main Controls (Prominent & Central):**  
   * **Container:** div with flex and justify-center, items-center.  
   * **"Shuffle" Toggle (Play/Pause):**  
     * **Tag:** \<button\> or \<a\> for interactivity.  
     * **Text:** "Shuffle" (when paused/ready to play) or an icon (e.g., pause symbol when playing).  
     * **Styling:**  
       * **Size:** Large, w-32 to w-48 (responsive). h-16 to h-24.  
       * **Shape:** rounded-full or rounded-xl (pill shape).  
       * **Background:** bg-accent1.  
       * **Text:** text-dark-background (or black if accent is light). font-bold, text-2xl to text-3xl.  
       * **Interaction:** cursor-pointer. Subtle hover:opacity-80 and active:scale-95.  
       * **State:** When playing, might show a different icon/text and have a subtle pulse animation or shadow-lg.  
   * **"Next" Button:**  
     * **Tag:** \<button\> or \<a\>.  
     * **Content:** Right-arrow icon (e.g., → or SVG) or "Next" text.  
     * **Styling:** w-12 to w-16, h-12 to h-16. rounded-full or rounded-lg. bg-accent1 or bg-accent2. text-dark-background or text-primary.  
     * **Placement:** To the right of the "Shuffle" toggle, with ml-4 or ml-8.  
     * **Interaction:** cursor-pointer. Subtle hover:opacity-80 and active:scale-95.  
2. **Search Field & Button:**  
   * **Container:** div with flex, items-center. mt-8 or mt-12 spacing.  
   * **Search Input Field:**  
     * **Tag:** \<input type="text"\>.  
     * **Placeholder:** "Search for a song, artist, or album...".  
     * **Styling:** w-full (or max-w-md on web). p-3 to p-4. rounded-lg. bg-gray-700 or bg-dark-background with a subtle border border-accent2. text-primary. focus:outline-none focus:ring-2 focus:ring-accent1.  
   * **Search Button:**  
     * **Tag:** \<button\>.  
     * **Content:** Magnifying glass icon (e.g., 🔍 or SVG).  
     * **Styling:** ml-3 or ml-4. p-3 to p-4. rounded-lg. bg-accent1. text-dark-background. cursor-pointer. hover:opacity-80.  
3. **Search Results List (Dynamic):**  
   * **Container:** div with overflow-y-auto. max-h-64 (or responsive height).  
   * **Styling:** mt-4. bg-dark-background or slightly lighter (\#222222).  
   * **Individual Result Item:**  
     * **Tag:** \<div\> or \<a\>.  
     * **Content:** Track title (font-bold, text-primary), Artist (text-sm, text-accent2), Album (text-xs, text-accent2).  
     * **Styling:** p-3 to p-4. border-b border-gray-800 (last item no border). cursor-pointer. hover:bg-gray-700 or hover:bg-accent1-opacity-10.  
     * **Interaction:** Clicking starts playback.

### **C. Bottom Section: Favourites List**

* **Container:** div taking remaining vertical space, overflow-y-auto.  
1. **"Favourites" Heading:**  
   * **Tag:** \<h2\> or \<div\>.  
   * **Text:** "Favourites".  
   * **Styling:** text-xl to text-2xl. font-bold. text-primary. my-6 (padding above/below).  
2. **Scrollable List of Favourites:**  
   * **Container:** div with overflow-y-auto. flex flex-col for vertical stacking.  
   * **Individual Favourite Item:**  
     * **Tag:** \<div\> or \<a\>.  
     * **Container:** div with flex items-center justify-between.  
     * **Content (Left):**  
       * Track Title (font-bold, text-primary).  
       * Artist (text-sm, text-accent2).  
       * Album (text-xs, text-accent2).  
     * **Content (Right):** "Untoggle Favourite" / Remove Button  
       * **Tag:** \<button\> or \<div\> with icon.  
       * **Icon:** ❤️ (filled heart) or ✕ (times/close icon).  
       * **Styling:** Small, text-red-500 or text-primary if ✕. cursor-pointer. hover:opacity-70.  
     * **Styling (Item):** p-3 to p-4. border-b border-gray-800. cursor-pointer. hover:bg-gray-700 or hover:bg-accent1-opacity-10.  
     * **Interaction:** Clicking anywhere on the item (except the remove button) starts playback. Clicking the icon/button removes from favorites.