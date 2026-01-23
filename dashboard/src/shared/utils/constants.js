
export const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';


export const CRITICALITY_LEVELS = {
  CRITICAL: 'critical',
  HIGH: 'high',
  MEDIUM: 'medium',
  LOW: 'low'
};


export const CRITICALITY_COLORS = {
  critical: '#dc2626',   
  high: '#ea580c',      
  medium: '#f59e0b',    
  low: '#22c55e'        
};


export const CRITICALITY_BG_CLASSES = {
  critical: 'bg-red-900/30 border-red-600 text-red-400',
  high: 'bg-orange-900/30 border-orange-600 text-orange-400',
  medium: 'bg-yellow-900/30 border-yellow-600 text-yellow-400',
  low: 'bg-green-900/30 border-green-600 text-green-400'
};


export const CRITICALITY_SCORE_RANGES = {
  critical: { min: 76, max: 100 },
  high: { min: 51, max: 75 },
  medium: { min: 26, max: 50 },
  low: { min: 0, max: 25 }
};


export const DEFAULT_PAGE_SIZE = 20;
export const PAGE_SIZE_OPTIONS = [10, 20, 50, 100];


export const DATE_FORMAT = {
  LONG: 'PPpp',           
  SHORT: 'PP',          
  TIME: 'p',            
  ISO: 'yyyy-MM-dd'       
};


export const CHART_COLORS = {
  critical: '#dc2626',
  high: '#ea580c',
  medium: '#f59e0b',
  low: '#22c55e',
  primary: '#3b82f6',     
  secondary: '#8b5cf6',  
  grid: '#374151'         
};


export const TIMELINE_DAYS_OPTIONS = [7, 14, 30];


export const SOURCE_CATEGORIES = [
  'forum',
  'marketplace',
  'news',
  'leak',
  'vulnerability',
  'ransomware',
  'other'
];


export const SCRAPE_INTERVALS = [
  { label: '30 minutes', value: '30 minutes' },
  { label: '1 hour', value: '1 hour' },
  { label: '4 hours', value: '4 hours' },
  { label: '12 hours', value: '12 hours' },
  { label: '24 hours', value: '24 hours' }
];


export const ERROR_CODES = {
  VALIDATION_ERROR: 'VALIDATION_ERROR',
  DATABASE_ERROR: 'DATABASE_ERROR',
  NOT_FOUND: 'NOT_FOUND',
  INVALID_ID: 'INVALID_ID',
  CREATE_FAILED: 'CREATE_FAILED',
  UPDATE_FAILED: 'UPDATE_FAILED',
  DELETE_FAILED: 'DELETE_FAILED'
};


export const FEATURE_TYPES = {
  BITCOIN: 'bitcoin_addresses',
  ETHEREUM: 'ethereum_addresses',
  MONERO: 'monero_addresses',
  ONION_URL: 'onion_urls',
  IP_ADDRESS: 'ip_addresses',
  EMAIL: 'emails',
  CVE: 'cves',
  KEYWORD: 'keywords'
};


export const FEATURE_DISPLAY_NAMES = {
  bitcoin_addresses: 'Bitcoin Addresses',
  ethereum_addresses: 'Ethereum Addresses',
  monero_addresses: 'Monero Addresses',
  onion_urls: 'Onion URLs',
  ip_addresses: 'IP Addresses',
  emails: 'Email Addresses',
  cves: 'CVEs',
  keywords: 'Keywords'
};


export const NOTIFICATION_TYPES = {
  SUCCESS: 'success',
  ERROR: 'error',
  WARNING: 'warning',
  INFO: 'info'
};


export const STORAGE_KEYS = {
  AUTH_TOKEN: 'cti_auth_token',
  USER_PREFERENCES: 'cti_user_preferences',
  THEME: 'cti_theme'
};