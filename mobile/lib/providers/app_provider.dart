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

  // Getters
  String? get userId => _userId;
  List<ClothingItem> get clothingItems => _clothingItems;
  List<Outfit> get outfits => _outfits;
  List<OutfitRecommendation> get recommendations => _recommendations;
  bool get isLoading => _isLoading;
  String? get error => _error;

  // Initialize
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

  // Load clothing items - sin loading para mostrar UI inmediatamente
  Future<void> loadClothingItems() async {
    try {
      _clothingItems = await _apiService.getClothingItems();
      _error = null;
      notifyListeners();
    } catch (e) {
      _error = e.toString();
      notifyListeners();
    }
  }

  // Load outfits
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

  // Get AI recommendations
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

  // Add clothing item
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

  // Upload image and create clothing item
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

  // Create outfit from recommendation
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
      // Actualizar la lista local
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
      // Actualizar la lista local
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
