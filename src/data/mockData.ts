
import { Novel, Chapter, SourceSite } from '../types/novel';

export const mockSites: SourceSite[] = [
  {
    id: 'qidian',
    name: 'Qidian',
    url: 'https://www.qidian.com',
    language: 'Chinese',
    icon: 'https://qidian.com/favicon.ico'
  },
  {
    id: 'naver',
    name: 'Naver Series',
    url: 'https://series.naver.com',
    language: 'Korean',
    icon: 'https://series.naver.com/favicon.ico'
  },
  {
    id: 'syosetu',
    name: 'Syosetu',
    url: 'https://syosetu.com',
    language: 'Japanese',
    icon: 'https://syosetu.com/favicon.ico'
  }
];

export const mockNovels: Novel[] = [
  {
    id: 'n1',
    title: 'The Supreme Magus',
    originalTitle: '至尊法师',
    cover: 'https://images.unsplash.com/photo-1531540758602-41bc4e95b2e4?ixlib=rb-1.2.1&auto=format&fit=crop&w=300&q=80',
    source: 'qidian',
    url: 'https://www.qidian.com/book/1234567',
    summary: 'A reincarnated mage seeks to master the ultimate magic in a world of cultivation and danger.',
    author: 'Jiang Xianxia',
    status: 'Ongoing',
    genres: ['Fantasy', 'Action', 'Adventure', 'Cultivation'],
    chaptersCount: 1243,
    lastUpdated: Date.now() - 2 * 24 * 60 * 60 * 1000,
    dateAdded: Date.now() - 30 * 24 * 60 * 60 * 1000,
  },
  {
    id: 'n2',
    title: 'Dragon Heart Chronicles',
    originalTitle: '드래곤 하트 연대기',
    cover: 'https://images.unsplash.com/photo-1604631806240-94d36a7a6626?ixlib=rb-1.2.1&auto=format&fit=crop&w=300&q=80',
    source: 'naver',
    url: 'https://series.naver.com/novel/detail.series?productNo=7654321',
    summary: 'A young hero with the soul of a dragon embarks on a journey to uncover his destiny.',
    author: 'Park Min-ho',
    status: 'Ongoing',
    genres: ['Fantasy', 'Adventure', 'Romance'],
    chaptersCount: 498,
    lastUpdated: Date.now() - 5 * 24 * 60 * 60 * 1000,
    dateAdded: Date.now() - 120 * 24 * 60 * 60 * 1000,
  },
  {
    id: 'n3',
    title: 'Silent Shadows',
    originalTitle: '静かな影',
    cover: 'https://images.unsplash.com/photo-1543002588-bfa74002ed7e?ixlib=rb-1.2.1&auto=format&fit=crop&w=300&q=80',
    source: 'syosetu',
    url: 'https://ncode.syosetu.com/n9876543',
    summary: 'In a world where shadows contain sentient beings, a young girl discovers she can communicate with them.',
    author: 'Tanaka Rei',
    status: 'Completed',
    genres: ['Mystery', 'Supernatural', 'Drama'],
    chaptersCount: 156,
    lastUpdated: Date.now() - 365 * 24 * 60 * 60 * 1000,
    dateAdded: Date.now() - 200 * 24 * 60 * 60 * 1000,
  },
  {
    id: 'n4',
    title: 'Immortal Sword Domain',
    originalTitle: '不朽剑域',
    cover: 'https://images.unsplash.com/photo-1601662528567-526cd06f6582?ixlib=rb-1.2.1&auto=format&fit=crop&w=300&q=80',
    source: 'qidian',
    url: 'https://www.qidian.com/book/2345678',
    summary: 'A legendary swordsman is reborn in the modern world and sets out to reclaim his heritage.',
    author: 'Li Tianxia',
    status: 'Ongoing',
    genres: ['Action', 'Fantasy', 'Martial Arts'],
    chaptersCount: 782,
    lastRead: {
      chapterId: 'c4-5',
      progress: 0.73,
      timestamp: Date.now() - 1 * 24 * 60 * 60 * 1000,
    },
    lastUpdated: Date.now() - 1 * 24 * 60 * 60 * 1000,
    dateAdded: Date.now() - 90 * 24 * 60 * 60 * 1000,
  },
  {
    id: 'n5',
    title: 'Virtual World Architect',
    originalTitle: '가상 세계 건축가',
    cover: 'https://images.unsplash.com/photo-1624437715218-272ca8e1bd31?ixlib=rb-1.2.1&auto=format&fit=crop&w=300&q=80',
    source: 'naver',
    url: 'https://series.naver.com/novel/detail.series?productNo=8765432',
    summary: 'A game developer gains the ability to manipulate reality like code and becomes a creator of worlds.',
    author: 'Kim Ji-woo',
    status: 'Ongoing',
    genres: ['Sci-Fi', 'LitRPG', 'Adventure'],
    chaptersCount: 324,
    lastUpdated: Date.now() - 3 * 24 * 60 * 60 * 1000,
    dateAdded: Date.now() - 60 * 24 * 60 * 60 * 1000,
  }
];

export const mockChapters: { [key: string]: Chapter[] } = {
  'n1': Array.from({ length: 10 }, (_, i) => ({
    id: `c1-${i+1}`,
    novelId: 'n1',
    number: i + 1,
    title: `Chapter ${i+1}: ${i === 0 ? 'The Beginning of Magic' : `Magical Journey Part ${i}`}`,
    originalTitle: `第${i+1}章: ${i === 0 ? '魔法的开始' : `魔法之旅 ${i}`}`,
    content: `<p>This is the translated content for chapter ${i+1} of The Supreme Magus.</p>
      <p>Leon looked up at the towering academy of magic, his heart racing with anticipation. After years of preparation, he was finally here—the prestigious Magus Academy, where only the most talented spellcasters were admitted.</p>
      <p>"Are you just going to stand there gawking all day?" a sharp voice cut through his reverie.</p>
      <p>Leon turned to see a tall girl with piercing blue eyes and silver hair that seemed to shimmer with magical energy. She wore the advanced student robes, marking her as at least a third-year.</p>
      <p>"I'm... I'm new here," Leon stammered, cursing his nervousness.</p>
      <p>The girl's expression softened slightly. "Obviously. I'm Selene, third-year elemental specialist. You'd better hurry—orientation starts in five minutes, and Master Thorne doesn't tolerate tardiness."</p>
      <p>Leon thanked her and rushed toward the grand entrance hall, his mind racing with possibilities. This was just the beginning of his journey to become the Supreme Magus, a title bestowed upon only the most powerful mage of each generation. Little did he know the challenges and dangers that awaited him beyond those ancient doors...</p>`,
    dateTranslated: Date.now() - (10 - i) * 24 * 60 * 60 * 1000,
    wordCount: 800 + Math.floor(Math.random() * 400)
  })),
  'n4': Array.from({ length: 8 }, (_, i) => ({
    id: `c4-${i+1}`,
    novelId: 'n4',
    number: i + 1,
    title: `Chapter ${i+1}: ${i === 0 ? 'The Sword Awakens' : `The Path of the Blade ${i}`}`,
    originalTitle: `第${i+1}章: ${i === 0 ? '剑醒' : `剑之道 ${i}`}`,
    content: `<p>This is the translated content for chapter ${i+1} of Immortal Sword Domain.</p>
      <p>The ancient sword hummed with power as Jiang Chen wrapped his fingers around its hilt. Centuries had passed since it had last been wielded, yet it responded to his touch as if it had been waiting specifically for him.</p>
      <p>"So the legends were true," whispered Master Wei, his old eyes wide with astonishment. "The Frost Emperor's Blade has chosen a new master."</p>
      <p>Jiang Chen could feel the cold energy of the sword flowing into his meridians, merging with his own qi. It was painful yet exhilarating, like being submerged in ice water while feeling more alive than ever before.</p>
      <p>"What does this mean?" Jiang Chen asked, though part of him already knew the answer.</p>
      <p>Master Wei stroked his long white beard thoughtfully. "It means, young man, that your destiny lies not in this small mountain village, but out there—in the vast sword domains where immortals clash and heavenly tribulations await the worthy."</p>
      <p>The old master's eyes glinted with a mixture of pride and concern. "The path of immortality through the sword is the most difficult of all cultivation roads. Many have tried; few have succeeded."</p>
      <p>Jiang Chen looked down at the gleaming blade, its surface etched with mysterious runes that seemed to shift and change as he watched. He thought of his humble beginnings, the tiny village that had been his whole world until now, and the parents he had lost to the demon beast invasion five years ago.</p>
      <p>"I'm ready," he said simply, and in that moment, the sword pulsed with blue light, as if acknowledging his resolve.</p>
      <p>What Jiang Chen didn't know was that far away, in the highest heavens of the cultivation world, certain immortals had felt the awakening of the Frost Emperor's Blade—and not all of them were pleased...</p>`,
    dateTranslated: Date.now() - (8 - i) * 24 * 60 * 60 * 1000,
    wordCount: 850 + Math.floor(Math.random() * 350)
  }))
};

// Function to get recently read novels
export const getRecentNovels = () => {
  return mockNovels
    .filter(novel => novel.lastRead)
    .sort((a, b) => (b.lastRead?.timestamp || 0) - (a.lastRead?.timestamp || 0))
    .slice(0, 4);
};

// Function to get recently updated novels
export const getRecentlyUpdatedNovels = () => {
  return mockNovels
    .sort((a, b) => b.lastUpdated - a.lastUpdated)
    .slice(0, 4);
};

// Function to search novels
export const searchNovels = (query: string) => {
  const lowerQuery = query.toLowerCase();
  return mockNovels.filter(novel => 
    novel.title.toLowerCase().includes(lowerQuery) ||
    novel.author?.toLowerCase().includes(lowerQuery) ||
    novel.summary.toLowerCase().includes(lowerQuery)
  );
};

// Function to get chapters for a novel
export const getChaptersForNovel = (novelId: string) => {
  return mockChapters[novelId] || [];
};

// Function to get a chapter
export const getChapter = (novelId: string, chapterNumber: number) => {
  const chapters = mockChapters[novelId] || [];
  return chapters.find(chapter => chapter.number === chapterNumber);
};

// Function to get a novel by ID
export const getNovelById = (id: string) => {
  return mockNovels.find(novel => novel.id === id);
};
