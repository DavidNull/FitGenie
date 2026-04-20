import 'dart:io';
import 'package:flutter/foundation.dart';
import '../models/clothing_item.dart';
import '../models/outfit.dart';
import '../models/outfit_recommendation.dart';
import '../services/api_service.dart';

class AppProvider extends ChangeNotifier {
  final ApiService _apiService = ApiService();
  
  String? _userId;
  List<ClothingItem> _clothingItems = [];
  List<Outfit> _outfits = [];
  List<OutfitRecommendation> _recommendations = [];
  bool _isLoading = false;
  String? _error;

  String? get userId => _userId;
  List<ClothingItem> get clothingItems => _clothingItems;
  List<Outfit> get outfits => _outfits;
  List<OutfitRecommendation> get recommendations => _recommendations;
  bool get isLoading => _isLoading;
  String? get error => _error;

  Future<void> initialize() async {
    _setLoading(true);
    try {
      final user = await _apiService.getCurrentUser();
      _userId = user['id'];
      await loadClothingItems();
      if (_userId != null) {
        await loadOutfits(_userId!);
      }
      _error = null;
    } catch (e) {
      _error = e.toString();
    }
    _setLoading(false);
  }

  final List<Map<String, dynamic>> _sampleImages = [
    {
      'path': 'assets/clothing/c1.png',
      'name': 'Camiseta Blanca',
      'category': 'Parte de arriba',
      'primaryColor': 'Blanco',
      'style': 'Casual',
      'season': ['Verano', 'Primavera'],
    },
    {
      'path': 'assets/clothing/c2.png',
      'name': 'Camiseta Negra',
      'category': 'Parte de arriba',
      'primaryColor': 'Negro',
      'style': 'Casual',
      'season': ['Verano', 'Otoño', 'Primavera'],
    },
    {
      'path': 'assets/clothing/c3.png',
      'name': 'Sudadera Gris',
      'category': 'Parte de arriba',
      'primaryColor': 'Gris',
      'style': 'Sport',
      'season': ['Invierno', 'Otoño'],
    },
    {
      'path': 'assets/clothing/p1.png',
      'name': 'Vaqueros Azules',
      'category': 'Parte de abajo',
      'primaryColor': 'Azul',
      'style': 'Casual',
      'season': ['Otoño', 'Invierno', 'Primavera'],
    },
    {
      'path': 'assets/clothing/p2.png',
      'name': 'Pantalón Negro',
      'category': 'Parte de abajo',
      'primaryColor': 'Negro',
      'style': 'Formal',
      'season': ['Otoño', 'Invierno', 'Primavera'],
    },
  ];

  Future<void> importSampleImages() async {
    _setLoading(true);
    try {
      for (final sample in _sampleImages) {
        final item = ClothingItem(
          id: '',
          userId: _userId ?? '',
          name: sample['name']!,
          category: sample['category']!,
          primaryColor: sample['primaryColor'],
          style: sample['style'],
          season: List<String>.from(sample['season'] ?? []),
          imageUrl: sample['path'],
        );
        await addClothingItem(item);
      }
      await loadClothingItems();
      _error = null;
    } catch (e) {
      _error = 'Error importing samples: $e';
    }
    _setLoading(false);
  }

  Future<void> loadClothingItems() async {
    _setLoading(true);
    try {
      _clothingItems = await _apiService.getClothingItems();
      _error = null;
      notifyListeners();
    } catch (e) {
      _error = 'Error al cargar prendas: $e';
      _clothingItems = [];
      notifyListeners();
    }
    _setLoading(false);
  }

  Future<void> loadOutfits(String userId) async {
    _setLoading(true);
    try {
      _outfits = await _apiService.getOutfits(userId);
      _error = null;
    } catch (e) {
      _error = e.toString();
    }
    _setLoading(false);
  }

  Future<void> getRecommendations({
    String? occasion,
    String? season,
    String? weather,
    String? style,
  }) async {
    if (_userId == null) return;
    
    _setLoading(true);
    try {
      final request = RecommendationRequest(
        occasion: occasion,
        season: season,
        weather: weather,
        style: style,
      );
      _recommendations = await _apiService.getRecommendations(
        _userId!,
        request: request,
      );
      _error = null;
    } catch (e) {
      _error = e.toString();
    }
    _setLoading(false);
  }

  Future<void> addClothingItem(ClothingItem item) async {
    _setLoading(true);
    try {
      final newItem = await _apiService.createClothingItem(item);
      _clothingItems.add(newItem);
      _error = null;
      notifyListeners();
    } catch (e) {
      _error = e.toString();
    }
    _setLoading(false);
  }

  Future<String?> uploadImage(File imageFile) async {
    _setLoading(true);
    try {
      final url = await _apiService.uploadImage(imageFile);
      _error = null;
      return url;
    } catch (e) {
      _error = e.toString();
      return null;
    } finally {
      _setLoading(false);
    }
  }

  Future<void> createOutfitFromRecommendation(OutfitRecommendation rec) async {
    if (_userId == null || rec.outfit == null) return;
    
    _setLoading(true);
    try {
      final outfit = await _apiService.createOutfit(rec.outfit!);
      _outfits.add(outfit);
      _error = null;
      notifyListeners();
    } catch (e) {
      _error = e.toString();
    }
    _setLoading(false);
  }

  Future<void> deleteClothingItem(String id) async {
    _setLoading(true);
    try {
      await _apiService.deleteClothingItem(id);
      _clothingItems.removeWhere((item) => item.id == id);
      _error = null;
      notifyListeners();
    } catch (e) {
      _error = e.toString();
    }
    _setLoading(false);
  }

  Future<void> updateClothingItem(ClothingItem item) async {
    _setLoading(true);
    try {
      final updatedItem = await _apiService.updateClothingItem(item);
      final index = _clothingItems.indexWhere((i) => i.id == item.id);
      if (index != -1) {
        _clothingItems[index] = updatedItem;
      }
      _error = null;
      notifyListeners();
    } catch (e) {
      _error = e.toString();
    }
    _setLoading(false);
  }

  Future<void> updateOutfit(Outfit outfit) async {
    _setLoading(true);
    try {
      final updatedOutfit = await _apiService.updateOutfit(outfit);
      final index = _outfits.indexWhere((o) => o.id == outfit.id);
      if (index != -1) {
        _outfits[index] = updatedOutfit;
      }
      _error = null;
      notifyListeners();
    } catch (e) {
      _error = e.toString();
    }
    _setLoading(false);
  }

  void _setLoading(bool value) {
    _isLoading = value;
    notifyListeners();
  }

  void clearError() {
    _error = null;
    notifyListeners();
  }
}
