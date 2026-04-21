import 'dart:convert';
import 'dart:io';
import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http;
import 'package:http_parser/http_parser.dart';
import '../models/clothing_item.dart';
import '../models/outfit.dart';
import '../models/outfit_recommendation.dart';

class ApiService {
  static String? _detectedHost;
  static bool _isDetecting = false;

  static String get apiHost {
    if (_detectedHost != null) return _detectedHost!;
    if (!_isDetecting) detectBackend();
    return _detectedHost ?? 'localhost';
  }

  static String get baseUrl => 'http://$apiHost:8080/api/v1';

  static Future<void> detectBackend() async {
    if (_isDetecting || _detectedHost != null) return;
    _isDetecting = true;

    final candidates = <String>[
      'localhost',
    ];

    // Android emulator
    if (!kIsWeb && Platform.isAndroid) {
      candidates.insert(0, '10.0.2.2');
    }

    // Linux/WSL: try common Docker bridge IPs
    if (!kIsWeb && Platform.isLinux) {
      candidates.addAll([
        '172.17.0.1',  // default docker0 bridge
        '172.21.48.1', // common WSL
        'host.docker.internal',
      ]);
    }

    for (final host in candidates) {
      try {
        final response = await http.get(
          Uri.parse('http://$host:8080/api/v1/users/me'),
          headers: {'X-Device-ID': 'auto-detect'},
        ).timeout(const Duration(seconds: 2));
        if (response.statusCode == 200) {
          _detectedHost = host;
          debugPrint('Backend detected at: $host');
          _isDetecting = false;
          return;
        }
      } catch (_) {
        debugPrint('Backend not at: $host');
      }
    }

    _detectedHost = candidates.first;
    _isDetecting = false;
  }
  
  String deviceId = 'flutter-test-device';
  String? userId;

  Map<String, String> get headers => {
    'Content-Type': 'application/json',
    'X-Device-ID': deviceId,
  };

  Future<Map<String, dynamic>> getCurrentUser() async {
    final response = await http.get(
      Uri.parse('${ApiService.baseUrl}/users/me'),
      headers: headers,
    ).timeout(const Duration(seconds: 5), onTimeout: () {
      throw Exception('Connection timeout - backend not reachable');
    });
    
    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      userId = data['id'];
      return data;
    }
    throw Exception('Failed to get user: ${response.statusCode}');
  }

  // ========== CLOTHING ITEMS ==========

  Future<List<ClothingItem>> getClothingItems() async {
    if (userId == null) {
      await getCurrentUser();
    }
    
    final response = await http.get(
      Uri.parse('${ApiService.baseUrl}/clothing?user_id=$userId'),
      headers: headers,
    ).timeout(const Duration(seconds: 5), onTimeout: () {
      throw Exception('Connection timeout - backend not reachable');
    });
    
    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      final items = data['items'] as List? ?? [];
      return items.map((e) => ClothingItem.fromJson(e)).toList();
    }
    throw Exception('Failed to get clothing items: ${response.statusCode}');
  }

  Future<ClothingItem> createClothingItem(ClothingItem item) async {
    final response = await http.post(
      Uri.parse('${ApiService.baseUrl}/clothing'),
      headers: headers,
      body: jsonEncode(item.toJson()),
    ).timeout(const Duration(seconds: 5));
    
    if (response.statusCode == 201) {
      return ClothingItem.fromJson(jsonDecode(response.body));
    }
    throw Exception('Failed to create clothing item: ${response.statusCode}');
  }

  Future<void> deleteClothingItem(String id) async {
    final response = await http.delete(
      Uri.parse('${ApiService.baseUrl}/clothing/$id'),
      headers: headers,
    ).timeout(const Duration(seconds: 5));
    
    if (response.statusCode != 200 && response.statusCode != 204) {
      throw Exception('Failed to delete clothing item: ${response.statusCode}');
    }
  }

  Future<ClothingItem> updateClothingItem(ClothingItem item) async {
    final response = await http.put(
      Uri.parse('${ApiService.baseUrl}/clothing/${item.id}'),
      headers: headers,
      body: jsonEncode(item.toJson()),
    ).timeout(const Duration(seconds: 5));
    
    if (response.statusCode == 200) {
      return ClothingItem.fromJson(jsonDecode(response.body));
    }
    throw Exception('Failed to update clothing item: ${response.statusCode}');
  }

  Future<List<Outfit>> getOutfits(String userId) async {
    final response = await http.get(
      Uri.parse('${ApiService.baseUrl}/users/$userId/outfits'),
      headers: headers,
    ).timeout(const Duration(seconds: 5));
    
    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      final outfits = data['outfits'] as List? ?? [];
      return outfits.map((e) => Outfit.fromJson(e)).toList();
    }
    throw Exception('Failed to get outfits: ${response.statusCode}');
  }

  Future<Outfit> createOutfit(Outfit outfit) async {
    if (userId == null) {
      await getCurrentUser();
    }
    
    final response = await http.post(
      Uri.parse('${ApiService.baseUrl}/outfits'),
      headers: headers,
      body: jsonEncode({
        'user_id': userId,
        'name': outfit.name,
        'description': outfit.description,
        'style': outfit.style,
        'occasion': outfit.occasion,
        'season': outfit.season,
        'weather': outfit.weather,
      }),
    ).timeout(const Duration(seconds: 5));
    
    if (response.statusCode == 201) {
      return Outfit.fromJson(jsonDecode(response.body));
    }
    throw Exception('Failed to create outfit: ${response.statusCode}');
  }

  Future<Outfit> updateOutfit(Outfit outfit) async {
    final response = await http.put(
      Uri.parse('${ApiService.baseUrl}/outfits/${outfit.id}'),
      headers: headers,
      body: jsonEncode({
        'name': outfit.name,
        'description': outfit.description,
        'style': outfit.style,
        'occasion': outfit.occasion,
        'season': outfit.season,
        'weather': outfit.weather,
        'rating': outfit.rating,
        'worn': outfit.worn,
        'favorite': outfit.favorite,
        'notes': outfit.notes,
      }),
    ).timeout(const Duration(seconds: 5));
    
    if (response.statusCode == 200) {
      return Outfit.fromJson(jsonDecode(response.body));
    }
    throw Exception('Failed to update outfit: ${response.statusCode}');
  }

  // ========== AI RECOMMENDATIONS ==========

  Future<List<OutfitRecommendation>> getRecommendations(
    String userId, {
    RecommendationRequest? request,
  }) async {
    final response = await http.post(
      Uri.parse('${ApiService.baseUrl}/users/$userId/outfits/recommendations'),
      headers: headers,
      body: jsonEncode(request?.toJson() ?? {}),
    ).timeout(const Duration(seconds: 10));
    
    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      final recommendations = data['recommendations'] as List? ?? [];
      return recommendations.map((e) => OutfitRecommendation.fromJson(e)).toList();
    }
    throw Exception('Failed to get recommendations: ${response.statusCode}');
  }

  // ========== UPLOAD ==========

  Future<String> uploadImage(File imageFile) async {
    final request = http.MultipartRequest(
      'POST',
      Uri.parse('${ApiService.baseUrl}/upload'),
    );
    request.headers['X-Device-ID'] = deviceId;
    final ext = imageFile.path.toLowerCase().split('.').last;
    final contentType = switch (ext) {
      'jpg' || 'jpeg' => MediaType('image', 'jpeg'),
      'png' => MediaType('image', 'png'),
      'webp' => MediaType('image', 'webp'),
      _ => MediaType('image', 'jpeg'),
    };
    
    request.files.add(await http.MultipartFile.fromPath(
      'image',
      imageFile.path,
      contentType: contentType,
    ));

    final response = await request.send().timeout(const Duration(seconds: 30));
    
    if (response.statusCode == 200) {
      final responseData = await response.stream.bytesToString();
      final data = jsonDecode(responseData);
      return data['url'] ?? '';
    }
    final errorData = await response.stream.bytesToString();
    throw Exception('Failed to upload image: ${response.statusCode} - $errorData');
  }

  // ========== COLOR THEORY ==========

  Future<Map<String, dynamic>> getColorSeasons() async {
    final response = await http.get(
      Uri.parse('${ApiService.baseUrl}/color-theory/seasons'),
      headers: headers,
    );
    
    if (response.statusCode == 200) {
      return jsonDecode(response.body);
    }
    throw Exception('Failed to get color seasons: ${response.statusCode}');
  }

  Future<Map<String, dynamic>> analyzeColorHarmony(List<String> colors) async {
    final response = await http.post(
      Uri.parse('${ApiService.baseUrl}/color-theory/analyze-harmony'),
      headers: headers,
      body: jsonEncode({'colors': colors}),
    );
    
    if (response.statusCode == 200) {
      return jsonDecode(response.body);
    }
    throw Exception('Failed to analyze harmony: ${response.statusCode}');
  }
}
