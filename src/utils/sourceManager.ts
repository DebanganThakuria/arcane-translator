import { SourceSite } from '../types/novel';
import { getAllSourceSites } from '../database/db';

class SourceManager {
  private sourceSites: SourceSite[] | null = null;
  private sourceLanguageMap: { [key: string]: string } = {};
  private loadingPromise: Promise<void> | null = null;

  async loadSources(): Promise<void> {
    // If already loaded, return immediately
    if (this.sourceSites !== null) {
      return;
    }

    // If already loading, return the existing promise
    if (this.loadingPromise) {
      return this.loadingPromise;
    }

    // Start loading
    this.loadingPromise = this.fetchSources();
    return this.loadingPromise;
  }

  private async fetchSources(): Promise<void> {
    try {
      const sources = await getAllSourceSites();
      this.sourceSites = sources;
      
      // Create the source ID to language flag mapping
      this.sourceLanguageMap = {};
      sources.forEach(source => {
        let flag = '🌐'; // Default globe
        switch (source.language.toLowerCase()) {
          case 'chinese':
            flag = '🇨🇳';
            break;
          case 'korean':
            flag = '🇰🇷';
            break;
          case 'japanese':
            flag = '🇯🇵';
            break;
          default:
            flag = '🌐';
        }
        this.sourceLanguageMap[source.id] = flag;
      });
    } catch (error) {
      console.error('Error loading source sites:', error);
      // Set empty arrays to avoid infinite retries
      this.sourceSites = [];
      this.sourceLanguageMap = {};
    } finally {
      this.loadingPromise = null;
    }
  }

  getLanguageFlag(sourceId: string): string {
    // First try the dynamic mapping from backend
    if (this.sourceLanguageMap[sourceId]) {
      return this.sourceLanguageMap[sourceId];
    }
    
    // Fallback to string matching for URLs or other formats
    const sourceLower = sourceId.toLowerCase();
    if (sourceLower.includes('chinese') || sourceLower.includes('cn') || sourceLower.includes('china')) {
      return '🇨🇳';
    }
    if (sourceLower.includes('korean') || sourceLower.includes('kr') || sourceLower.includes('korea')) {
      return '🇰🇷';
    }
    if (sourceLower.includes('japanese') || sourceLower.includes('jp') || sourceLower.includes('japan')) {
      return '🇯🇵';
    }
    
    // Default to globe for unknown sources
    return '🌐';
  }

  getSourceSites(): SourceSite[] {
    return this.sourceSites || [];
  }

  getSourceById(id: string): SourceSite | undefined {
    return this.sourceSites?.find(source => source.id === id);
  }

  isLoaded(): boolean {
    return this.sourceSites !== null;
  }
}

// Export a singleton instance
export const sourceManager = new SourceManager(); 