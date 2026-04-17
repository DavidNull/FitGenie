class ClothingItem {
  final String id;
  final String userId;
  final String name;
  final String category;
  final String? brand;
  final String? size;
  final String? primaryColor;
  final String? secondaryColor;
  final String? material;
  final String? style;
  final List<String> season;
  final List<String> occasion;
  final String? imageUrl;
  final String? notes;
  final DateTime? createdAt;
  final DateTime? updatedAt;
  final bool isLocalAsset;

  ClothingItem({
    required this.id,
    required this.userId,
    required this.name,
    required this.category,
    this.brand,
    this.size,
    this.primaryColor,
    this.secondaryColor,
    this.material,
    this.style,
    this.season = const [],
    this.occasion = const [],
    this.imageUrl,
    this.notes,
    this.createdAt,
    this.updatedAt,
    this.isLocalAsset = false,
  });

  factory ClothingItem.fromJson(Map<String, dynamic> json) {
    return ClothingItem(
      id: json['id'] ?? '',
      userId: json['user_id'] ?? '',
      name: json['name'] ?? '',
      category: json['category'] ?? '',
      brand: json['brand'],
      size: json['size'],
      primaryColor: json['primary_color'],
      secondaryColor: json['secondary_color'],
      material: json['material'],
      style: json['style'],
      season: List<String>.from(json['season'] ?? []),
      occasion: List<String>.from(json['occasion'] ?? []),
      imageUrl: json['image_url'],
      notes: json['notes'],
      createdAt: json['created_at'] != null 
          ? DateTime.tryParse(json['created_at']) 
          : null,
      updatedAt: json['updated_at'] != null 
          ? DateTime.tryParse(json['updated_at']) 
          : null,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'user_id': userId,
      'name': name,
      'category': category,
      'brand': brand,
      'size': size,
      'primary_color': primaryColor,
      'secondary_color': secondaryColor,
      'material': material,
      'style': style,
      'season': season,
      'occasion': occasion,
      'image_url': imageUrl,
      'notes': notes,
      if (createdAt != null) 'created_at': createdAt!.toIso8601String(),
      if (updatedAt != null) 'updated_at': updatedAt!.toIso8601String(),
    };
  }
}
