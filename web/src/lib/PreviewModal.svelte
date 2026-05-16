<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  const dispatch = createEventDispatcher();

  export let file: any;
  export let isOpen: boolean;

  const API_URL = 'http://localhost:8080';

  function close() {
    dispatch('close');
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') close();
  }
</script>

<svelte:window on:keydown={handleKeydown}/>

{#if isOpen && file}
  <div class="fixed inset-0 z-[100] flex items-center justify-center p-4 sm:p-8 bg-black/95 backdrop-blur-sm" on:click={close}>
    <div class="relative max-w-5xl w-full max-h-full flex flex-col items-center gap-4" on:click|stopPropagation>
      
      <!-- Close Button -->
      <button 
        on:click={close}
        class="absolute -top-12 right-0 text-white/70 hover:text-white transition-colors"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>

      <!-- Content -->
      <div class="w-full h-full flex items-center justify-center overflow-hidden rounded-lg shadow-2xl bg-gray-900/50">
        {#if file.extension.toLowerCase().match(/\.(jpg|jpeg|png|gif|webp|bmp)$/)}
          <img 
            src={`${API_URL}/media/${file.path}`} 
            alt={file.name}
            class="max-w-full max-h-[80vh] object-contain"
          />
        {:else if file.extension.toLowerCase().match(/\.(mp4|webm|mkv|avi)$/)}
          <video 
            src={`${API_URL}/media/${file.path}`} 
            controls 
            autoplay
            class="max-w-full max-h-[80vh]"
          ></video>
        {:else}
          <div class="p-20 text-center space-y-4">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-20 w-20 mx-auto text-gray-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <p class="text-gray-400">Preview not available for this file type</p>
          </div>
        {/if}
      </div>

      <!-- Footer Info -->
      <div class="w-full text-center space-y-1">
        <h2 class="text-xl font-semibold text-white truncate">{file.name}</h2>
        <div class="flex items-center justify-center gap-4 text-sm text-gray-400">
          <span>{file.size}</span>
          <span>•</span>
          <span>{file.date_modified}</span>
          <span>•</span>
          <span class="uppercase">{file.extension.replace('.','')}</span>
        </div>
      </div>

    </div>
  </div>
{/if}
